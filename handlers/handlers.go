package handlers

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Разрешаем все источники (для локальной разработки)
		return true
	},
}

var (
	clients   = make(map[*websocket.Conn]bool) // Подключённые клиенты
	clientsMu sync.Mutex                       // Мьютекс для потокобезопасности
)

// HandleConnections обрабатывает входящие WebSocket-соединения.
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Регистрируем нового клиента
	clientsMu.Lock()
	clients[ws] = true
	clientsMu.Unlock()

	log.Println("Новое соединение установлено")

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("Ошибка чтения сообщения:", err)
			clientsMu.Lock()
			delete(clients, ws)
			clientsMu.Unlock()
			return
		}

		log.Printf("Получено сообщение: %s", message)

		// Пересылаем сообщение всем клиентам
		clientsMu.Lock()
		for client := range clients {
			if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("Ошибка отправки сообщения:", err)
				client.Close()
				delete(clients, client)
			}
		}
		clientsMu.Unlock()
	}
}
