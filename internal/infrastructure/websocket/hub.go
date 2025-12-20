package websocket

import (
	"encoding/json"
	"sync"
	"time"
)

type Hub struct {
	clients map[uint]map[*Client]bool
	rooms map[uint]map[*Client]bool

	Broadcast   chan *Message
	Register    chan *Client
	Unregister  chan *Client
	mu          sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]map[*Client]bool),
		rooms:      make(map[uint]map[*Client]bool),
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

	presenceBytes := h.buildPresenceBytes(c.roomID, c.userID, true)

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
		}
	}

	if set, ok := h.rooms[c.roomID]; ok {
		delete(set, c)
		if len(set) == 0 {
			delete(h.rooms, c.roomID)
		}
	}

	presenceBytes := h.buildPresenceBytes(c.roomID, c.userID, false)
	targets := h.snapshotRoomLocked(c.roomID)

	h.mu.Unlock()

	c.CloseConnection()

	h.sendToClients(targets, presenceBytes)
}

func (h *Hub) handleBroadcast(msg *Message) {
	if msg.Type != MessageTypeChat {
		return
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return
	}

	h.mu.RLock()
	targets := h.snapshotRoomRLocked(msg.RoomID)
	h.mu.RUnlock()

	h.sendToClients(targets, payload)
}

func (h *Hub) IsUserOnline(userID uint) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	set, ok := h.clients[userID]
	return ok && len(set) > 0
}


func (h *Hub) buildPresenceBytes(roomID, userID uint, online bool) []byte {
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
