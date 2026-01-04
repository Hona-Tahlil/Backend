package websocket

import (
	"encoding/json"
	"sync"
	"time"
)

type Hub struct {
	clients map[uint]map[*Client]bool
	rooms map[uint]map[*Client]bool
	lastSeen map[uint]time.Time

	Broadcast   chan *Message
	Register    chan *Client
	Unregister  chan *Client
	mu          sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]map[*Client]bool),
		rooms:      make(map[uint]map[*Client]bool),
		lastSeen:   make(map[uint]time.Time),
		Broadcast:  make(chan *Message, 256),
		Register:   make(chan *Client, 256),
		Unregister: make(chan *Client, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.handleRegister(c)
		case c := <-h.Unregister:
			h.handleUnregister(c)
		case msg := <-h.Broadcast:
			h.handleBroadcast(msg)
		}
	}
}

func (h *Hub) handleRegister(c *Client) {
	h.mu.Lock()

	if _, ok := h.clients[c.userID]; !ok {
		h.clients[c.userID] = make(map[*Client]bool)
	}
	h.clients[c.userID][c] = true

	if _, ok := h.rooms[c.roomID]; !ok {
		h.rooms[c.roomID] = make(map[*Client]bool)
	}
	h.rooms[c.roomID][c] = true

	presenceBytes := h.buildPresenceBytes(c.roomID, c.userID, true, nil)

	targets := h.snapshotRoomLocked(c.roomID)

	h.mu.Unlock()

	h.sendToClients(targets, presenceBytes)
}

func (h *Hub) handleUnregister(c *Client) {
	h.mu.Lock()

	if set, ok := h.clients[c.userID]; ok {
		delete(set, c)
		if len(set) == 0 {
			delete(h.clients, c.userID)
			h.lastSeen[c.userID] = time.Now().UTC()
		}
	}
	isOnlineGlobal := false
	if set, ok := h.clients[c.userID]; ok && len(set) > 0 {
		isOnlineGlobal = true
	}

	if set, ok := h.rooms[c.roomID]; ok {
		delete(set, c)
		if len(set) == 0 {
			delete(h.rooms, c.roomID)
		}
	}
	stillInRoom := h.isUserOnlineInRoomLocked(c.roomID, c.userID)

	var presenceBytes []byte
	var targets []*Client
	if !stillInRoom {
		var lastSeen *time.Time
		if !isOnlineGlobal {
			ts := h.lastSeen[c.userID]
			lastSeen = &ts
		}
		presenceBytes = h.buildPresenceBytes(c.roomID, c.userID, false, lastSeen)
		targets = h.snapshotRoomLocked(c.roomID)
	}

	h.mu.Unlock()

	c.CloseConnection()

	if len(presenceBytes) != 0 {
		h.sendToClients(targets, presenceBytes)
	}
}

func (h *Hub) handleBroadcast(msg *Message) {
	h.mu.RLock()
	targets := h.snapshotRoomRLocked(msg.RoomID)
	h.mu.RUnlock()

	switch msg.Type {
	case MessageTypeChat:
		h.sendChatToTargets(targets, msg)
	case MessageTypeRead:
		fallthrough
	case MessageTypeEdit:
		fallthrough
	case MessageTypeDelete:
		fallthrough
	case MessageTypeReaction:
		payload, err := json.Marshal(msg)
		if err != nil {
			return
		}
		h.sendToClients(targets, payload)
	default:
		return
	}
}

func (h *Hub) IsUserOnline(userID uint) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	set, ok := h.clients[userID]
	return ok && len(set) > 0
}

func (h *Hub) LastSeen(userID uint) *time.Time {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if ts, ok := h.lastSeen[userID]; ok {
		copy := ts
		return &copy
	}
	return nil
}

func (h *Hub) buildPresenceBytes(roomID, userID uint, online bool, lastSeen *time.Time) []byte {
	wire := struct {
		Type    string          `json:"type"`
		RoomID  uint            `json:"room_id"`
		Content PresencePayload `json:"content"`
	}{
		Type:   MessageTypePresence,
		RoomID: roomID,
		Content: PresencePayload{
			UserID:    userID,
			IsOnline:  online,
			LastSeen:  lastSeen,
			Timestamp: time.Now(),
		},
	}
	b, _ := json.Marshal(wire)
	return b
}

func (h *Hub) snapshotRoomLocked(roomID uint) []*Client {
	set, ok := h.rooms[roomID]
	if !ok || len(set) == 0 {
		return nil
	}
	out := make([]*Client, 0, len(set))
	for c := range set {
		out = append(out, c)
	}
	return out
}

func (h *Hub) snapshotRoomRLocked(roomID uint) []*Client {
	set, ok := h.rooms[roomID]
	if !ok || len(set) == 0 {
		return nil
	}
	out := make([]*Client, 0, len(set))
	for c := range set {
		out = append(out, c)
	}
	return out
}

func (h *Hub) isUserOnlineInRoomLocked(roomID, userID uint) bool {
	set, ok := h.rooms[roomID]
	if !ok || len(set) == 0 {
		return false
	}
	for c := range set {
		if c.userID == userID {
			return true
		}
	}
	return false
}

func (h *Hub) sendToClients(targets []*Client, payload []byte) {
	if len(targets) == 0 {
		return
	}
	for _, c := range targets {
		select {
		case <-c.Done():
			select {
			case h.Unregister <- c:
			default:
			}
			continue
		default:
		}

		select {
		case c.send <- payload:
		default:
			select {
			case h.Unregister <- c:
			default:
			}
		}
	}
}

func (h *Hub) sendChatToTargets(targets []*Client, msg *Message) {
	if len(targets) == 0 || msg == nil {
		return
	}

	var payload ChatPayload
	if err := json.Unmarshal(msg.Content, &payload); err != nil {
		wire, err := json.Marshal(msg)
		if err != nil {
			return
		}
		h.sendToClients(targets, wire)
		return
	}

	for _, c := range targets {
		outPayload := payload
		outPayload.IsMine = c.userID == msg.SenderID
		outPayload.IsUnread = !outPayload.IsMine

		wire := struct {
			Type      string      `json:"type"`
			RoomID    uint        `json:"room_id"`
			SenderID  uint        `json:"sender_id,omitempty"`
			MessageID uint        `json:"message_id,omitempty"`
			Timestamp time.Time   `json:"timestamp"`
			Content   ChatPayload `json:"content,omitempty"`
		}{
			Type:      msg.Type,
			RoomID:    msg.RoomID,
			SenderID:  msg.SenderID,
			MessageID: msg.MessageID,
			Timestamp: msg.Timestamp,
			Content:   outPayload,
		}

		b, err := json.Marshal(wire)
		if err != nil {
			continue
		}
		h.sendToClients([]*Client{c}, b)
	}
}
