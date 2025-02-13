package handlers

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// HandleConnections обрабатывает входящие WebSocket-соединения
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Обновляем соединение до WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Канал для асинхронной отправки сообщений
	msgChan := make(chan []byte)
	defer close(msgChan)

	// Горутина для отправки сообщений
	go func() {
		for message := range msgChan {
			if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("Ошибка отправки сообщения:", err)
				break
			}
		}
	}()

	// Настройка обработчика Ping
	ws.SetPingHandler(func(appData string) error {
		log.Println("Получен Ping")
		return ws.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(10*time.Second))
	})

	// Настройка обработчика Pong
	ws.SetPongHandler(func(appData string) error {
		log.Println("Получен Pong")
		return nil
	})

	// Таймер для отправки Ping
	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()

	// Основной цикл обработки сообщений
	for {
		select {
		case <-pingTicker.C:
			// Отправляем Ping клиенту
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				log.Println("Ошибка отправки Ping:", err)
				return
			}

		default:
			// Чтение сообщения от клиента
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("Ошибка чтения сообщения:", err)
				return
			}

			// Логируем полученное сообщение
			log.Printf("Received: %s", message)

			// Отправляем сообщение в канал для асинхронной обработки
			msgChan <- message
		}
	}
}
