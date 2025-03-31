package main

import (
	"net/http"
	"web-server/server"
)

func main() {
	// Заглушка для favicon.ico
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "favicon.ico")
	})

	// Запускаем сервер
	server.StartServer(":9000")
}
