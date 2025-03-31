package server

import (
	"log"
	"net/http"
	"web-server/handlers"
)

// StartServer запускает сервер на указанном порту.
func StartServer(port string) {
	// Регистрируем обработчик WebSocket
	http.HandleFunc("/ws", handlers.HandleConnections)

	log.Printf("Сервер запущен на %s\n", port)

	// Запускаем HTTP сервер
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
