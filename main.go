package main

import (
	"log"
	"net/http"
	"web-server/server"
)

func main() {
	// Заглушка для favicon.ico
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "favicon.ico")
	})

	// Запуск сервера на порту 8096
	err := server.StartServer(":8096")
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
