package client

import (
	"github.com/gorilla/websocket"
	"io"
)

// Client представляет WebSocket клиента.
type Client struct {
	ID     int
	WS     *websocket.Conn
	Server *Server
	Ch     chan *Message
	DoneCh chan bool
}

// Listen начинает прослушивание сообщений от клиента.
func (client *Client) Listen() {
	go client.listenWrite()
	client.listenRead()
}

func (client *Client) listenWrite() {
	for {
		select {
		case msg := <-client.Ch:
			err := client.WS.WriteJSON(msg)
			if err != nil {
				client.Server.Err(err)
				client.DoneCh <- true
				return
			}
		case <-client.DoneCh:
			client.Server.Del(client)
			client.DoneCh <- true
			return
		}
	}
}

func (client *Client) listenRead() {
	for {
		var msg Message
		err := client.WS.ReadJSON(&msg)
		if err == io.EOF {
			client.DoneCh <- true
			return
		} else if err != nil {
			client.Server.Err(err)
			return
		}
		client.Server.SendAll(&msg) // Отправляем сообщение всем клиентам
	}
}
