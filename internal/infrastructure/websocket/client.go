package websocket

import (
	"bytes"
	"encoding/json"
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/chat"
	"hona/backend/internal/application/usecase"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	websocketSetting *bootstrap.WebsocketSetting
	Hub              *Hub

	conn *websocket.Conn
	send chan []byte

	roomID uint
	userID uint

	mu        sync.RWMutex
	done      chan struct{}
	closeOnce sync.Once
	isReady   bool

	chatService usecase.ChatService
}

func NewClient(
	hub *Hub,
	conn any,
	roomID, userID uint,
	websocketSetting *bootstrap.WebsocketSetting,
	chatService usecase.ChatService,
) *Client {
	wsConn, _ := conn.(*websocket.Conn)

	return &Client{
		websocketSetting: websocketSetting,
		Hub:              hub,
		conn:             wsConn,
		send:             make(chan []byte, websocketSetting.MessageBufferSize),
		roomID:           roomID,
		userID:           userID,
		done:             make(chan struct{}),
		chatService:      chatService,
	}
}

func (c *Client) ReadPump() error {
	// مهم: وقتی ReadPump تمام شد، کلاینت را از هاب Unregister کن
	defer func() {
		// Unregister باعث CloseConnection هم می‌شود
		c.Hub.Unregister <- c
	}()

	c.mu.Lock()
	c.isReady = true
	c.mu.Unlock()

	c.conn.SetReadLimit(int64(c.websocketSetting.MaxMessageSize))
	_ = c.conn.SetReadDeadline(time.Now().Add(c.websocketSetting.ReadTimeout))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(c.websocketSetting.ReadTimeout))
		return nil
	})

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			return err
		}

		var msg Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}

		msg.Client = c
		msg.RoomID = c.roomID
		msg.SenderID = c.userID
		msg.Timestamp = time.Now()

		switch msg.Type {
		case MessageTypeChat:
			c.processAndSaveChatMessage(&msg)
			// اگر ذخیره نشد، broadcast نکن
			if msg.MessageID == 0 {
				continue
			}
			c.Hub.Broadcast <- &msg
		default:
			// پیام‌های ناشناس را ignore کن
			continue
		}
	}
}

func (c *Client) WritePump() error {
	ticker := time.NewTicker(c.websocketSetting.PingPeriod)
	defer func() {
		ticker.Stop()
		// اگر WritePump تمام شد، Unregister کن (اگر قبلاً نشده)
		select {
		case c.Hub.Unregister <- c:
		default:
			// اگر کانال پر بود/هاب بسته بود، فقط کانکشن را ببند
			c.CloseConnection()
		}
	}()

	for {
		select {
		case <-c.done:
			return nil

		case payload, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(c.websocketSetting.WriteTimeout))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closed by server"))
				return nil
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return err
			}

			_, _ = w.Write(payload)

			// batch: پیام‌های در صف را هم پشت سر هم بفرست
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(bytes.TrimSpace([]byte{'\n'}))
				_, _ = w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return err
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(c.websocketSetting.WriteTimeout))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return err
			}
		}
	}
}

func (c *Client) CloseConnection() {
	c.closeOnce.Do(func() {
		close(c.done)

		// close(send) را اینجا انجام می‌دهیم چون WritePump روی send می‌خواند
		// اگر send قبلاً بسته شده باشد panic می‌دهد، پس با recover یا once کنترل می‌کنیم
		defer func() { _ = recover() }()
		close(c.send)

		_ = c.conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closing connection"),
		)
		_ = c.conn.Close()
	})
}

func (c *Client) processAndSaveChatMessage(msg *Message) {
	var content string
	if err := json.Unmarshal(msg.Content, &content); err != nil {
		return
	}
	req := chat.SaveMessageRequest{
		RoomID:  c.roomID,
		SenderID:  c.userID,
		Content: content,
		ReplyToMessageID: nil,
	}
	saved, err := c.chatService.SaveMessage(req)
	if err != nil {
		return
	}

	msg.MessageID = saved.ID
}

func (c *Client) IsReady() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isReady
}

func (c *Client) UserID() uint { return c.userID }
func (c *Client) RoomID() uint { return c.roomID }
func (c *Client) Done() <-chan struct{} {
	return c.done
}
