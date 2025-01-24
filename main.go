package main

import (
	"fmt"
	"go-chat-server/server"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", server.HandleConnections)
	go server.HandleMessages()

	fmt.Println("Сервер запущен на порту 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска:", err)
	}
}
