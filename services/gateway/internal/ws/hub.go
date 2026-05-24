package ws

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type gameSession struct {
	count     int
	sessionID string
}

type Hub struct {
	mu           sync.Mutex
	clients      map[*websocket.Conn]struct{}
	gameSessions map[string]gameSession
}

func NewHub() *Hub {
	return &Hub{
		clients:      make(map[*websocket.Conn]struct{}),
		gameSessions: make(map[string]gameSession),
	}
}

func (h *Hub) Broadcast(payload any) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for conn := range h.clients {
		if err := conn.WriteJSON(payload); err != nil {
			log.Printf("write websocket message: %v", err)
			_ = conn.Close()
			delete(h.clients, conn)
		}
	}
}

func (h *Hub) add(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[conn] = struct{}{}
}

func (h *Hub) remove(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	_ = conn.Close()
	delete(h.clients, conn)
}

func (h *Hub) ActivateGame(gameID, sessionID string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	session := h.gameSessions[gameID]
	if session.count == 0 {
		session.sessionID = sessionID
	}
	session.count++
	h.gameSessions[gameID] = session

	return session.count == 1
}

func (h *Hub) DeactivateGame(gameID string) (string, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	session, ok := h.gameSessions[gameID]
	if !ok {
		return "", false
	}

	session.count--
	if session.count > 0 {
		h.gameSessions[gameID] = session
		return session.sessionID, false
	}

	delete(h.gameSessions, gameID)
	return session.sessionID, true
}
