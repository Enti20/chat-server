package server

import (
	"log"
	"net/http"
	"web-server/handlers"
)

// StartServer запускает сервер на указанном порту.
func StartServer(port string) error {
	// Регистрируем обработчик WebSocket
	http.HandleFunc("/ws", handlers.HandleConnections)

	log.Printf("Сервер запущен на %s\n", port)
	return http.ListenAndServe(port, nil)
}
