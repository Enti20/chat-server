package main

import (
	"log"
	"web-server/server"
)

func main() {
	// Запуск сервера на порту 8080
	err := server.StartServer(":8080")
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
