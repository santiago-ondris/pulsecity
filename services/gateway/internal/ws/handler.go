package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

func (h *Hub) ServeWebSocket(w http.ResponseWriter, r *http.Request, onConnect func(*websocket.Conn) error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade websocket: %v", err)
		return
	}

	h.add(conn)
	defer h.remove(conn)

	if onConnect != nil {
		if err := onConnect(conn); err != nil {
			log.Printf("initialize websocket: %v", err)
			return
		}
	}

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}
