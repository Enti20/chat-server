package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var (
	clients         = make(map[*websocket.Conn]int)
	clientsMu       sync.Mutex
	clientIDCounter int
)

// Message представляет сообщение чата.
type Message struct {
	ClientID int    `json:"client_id"`
	Author   string `json:"author"`
	Body     string `json:"body"`
}

// HandleConnections обрабатывает входящие WebSocket-соединения.
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка при установке соединения:", err)
		return
	}
	defer ws.Close()

	clientsMu.Lock()
	clientIDCounter++
	clientID := clientIDCounter
	clients[ws] = clientID
	clientsMu.Unlock()

	log.Printf("Клиент %d подключился", clientID)

	broadcastMessage(Message{
		ClientID: clientID,
		Author:   "Система",
		Body:     fmt.Sprintf("Клиент %d подключился", clientID),
	})

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Клиент %d отключился", clientID)
			clientsMu.Lock()
			delete(clients, ws)
			clientsMu.Unlock()

			broadcastMessage(Message{
				ClientID: clientID,
				Author:   "Система",
				Body:     fmt.Sprintf("Клиент %d отключился", clientID),
			})

			return
		}

		log.Printf("Сообщение от клиента %d: %s", clientID, message)

		msg := Message{
			ClientID: clientID,
			Author:   fmt.Sprintf("Клиент %d", clientID),
			Body:     string(message),
		}

		broadcastMessage(msg)
	}
}

func broadcastMessage(msg Message) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for client := range clients {
		if err := client.WriteJSON(msg); err != nil {
			log.Println("Ошибка отправки сообщения:", err)
			client.Close()
			delete(clients, client)
		}
	}
}
