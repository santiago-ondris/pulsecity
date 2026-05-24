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

func (h *Hub) ServeWebSocket(w http.ResponseWriter, r *http.Request, onConnect func(*websocket.Conn) (func(), error)) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade websocket: %v", err)
		return
	}

	var onDisconnect func()
	h.add(conn)
	defer func() {
		if onDisconnect != nil {
			onDisconnect()
		}
		h.remove(conn)
	}()

	if onConnect != nil {
		disconnect, err := onConnect(conn)
		if err != nil {
			log.Printf("initialize websocket: %v", err)
			return
		}
		onDisconnect = disconnect
	}

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}
