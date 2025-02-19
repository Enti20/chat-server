package server

import (
	"log"
	"net/http"
	"web-server/client"
)

// StartServer запускает сервер на указанном порту.
func StartServer(port string) {
	// Создаем новый сервер WebSocket
	wsServer := &client.Server{
		Pattern:   "/ws",
		Messages:  []*client.Message{},
		Clients:   make(map[int]*client.Client),
		AddCh:     make(chan *client.Client),
		DelCh:     make(chan *client.Client),
		SendAllCh: make(chan *client.Message),
		DoneCh:    make(chan bool),
		ErrCh:     make(chan error),
	}

	// Запускаем сервер
	go wsServer.Listen()

	log.Printf("Сервер запущен на %s\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
