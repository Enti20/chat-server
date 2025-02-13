package server

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"web-server/handlers"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// StartServer запускает сервер на указанном порту
func StartServer(port string) error {
	// Обработчик для корневого пути
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("WebSocket server is running!\n\n This is best of ze best in pet-project"))
	})

	// Регистрируем обработчик WebSocket
	http.HandleFunc("/ws", handlers.HandleConnections)

	log.Printf("Сервер запущен на %s\n", port)
	return http.ListenAndServe(port, nil)
}
