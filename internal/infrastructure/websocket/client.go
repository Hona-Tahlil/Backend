package websocket

import (
	"bytes"
	"encoding/json"
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/chat"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/enums"
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
	defer func() {
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
			if msg.MessageID == 0 {
				continue
			}
			c.Hub.Broadcast <- &msg
		case MessageTypeRead:
			c.processAndBroadcastReadReceipt(&msg)
		default:
			continue
		}
	}
}

func (c *Client) WritePump() error {
	ticker := time.NewTicker(c.websocketSetting.PingPeriod)
	defer func() {
		ticker.Stop()
		select {
		case c.Hub.Unregister <- c:
		default:
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
	var payload ChatPayload
	if err := json.Unmarshal(msg.Content, &payload); err != nil {
		var legacyContent string
		if err := json.Unmarshal(msg.Content, &legacyContent); err != nil {
			return
		}
		payload.MessageType = "TEXT"
		payload.Text = legacyContent
	}
	if payload.MessageType == "" {
		payload.MessageType = "TEXT"
	}
	req := chat.SaveMessageRequest{
		RoomID:      c.roomID,
		SenderID:    c.userID,
		Content:     payload.Text,
		MessageType: enums.ChatMessageType(payload.MessageType),
		MediaBase64: payload.ImageBase64,
		MediaMime:   payload.ImageMime,
		ReplyToMessageID: payload.ReplyToMessageID,
	}
	saved, err := c.chatService.SaveMessage(req)
	if err != nil {
		return
	}

	responsePayload := ChatPayload{
		MessageType: string(saved.MessageType),
		Text:        saved.Content,
		ReplyToMessageID: payload.ReplyToMessageID,
	}
	if saved.MediaURL != nil {
		responsePayload.ImageURL = *saved.MediaURL
	}
	encoded, err := json.Marshal(responsePayload)
	if err != nil {
		return
	}

	msg.MessageID = saved.ID
	msg.SenderID = c.userID
	msg.Content = encoded
}

func (c *Client) processAndBroadcastReadReceipt(msg *Message) {
	var payload ReadPayload
	if err := json.Unmarshal(msg.Content, &payload); err != nil {
		return
	}
	if payload.LastReadMessageID == 0 {
		return
	}

	req := chat.MarkRoomReadRequest{
		RoomID:            c.roomID,
		SenderID:          c.userID,
		LastReadMessageID: payload.LastReadMessageID,
	}
	if err := c.chatService.MarkAsRead(req); err != nil {
		return
	}

	payload.ReaderID = c.userID
	payload.Timestamp = time.Now()
	encoded, err := json.Marshal(payload)
	if err != nil {
		return
	}

	msg.RoomID = c.roomID
	msg.SenderID = c.userID
	msg.Timestamp = payload.Timestamp
	msg.Content = encoded
	c.Hub.Broadcast <- msg
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
