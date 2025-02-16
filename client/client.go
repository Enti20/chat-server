package client

import (
	"github.com/gorilla/websocket"
	"io"
)

// Client represents a WebSocket client.
type Client struct {
	ID     int
	WS     *websocket.Conn
	Server *Server
	Ch     chan *Message
	DoneCh chan bool
}

// Listen starts listening for messages from the client.
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
		select {
		case <-client.DoneCh:
			client.Server.Del(client)
			client.DoneCh <- true
			return
		default:
			var msg Message
			err := client.WS.ReadJSON(&msg)
			if err == io.EOF {
				client.DoneCh <- true
			} else if err != nil {
				client.Server.Err(err)
			} else {
				client.Server.SendAll(&msg)
			}
		}
	}
}

// Message represents a chat message.
type Message struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}

func (msg *Message) String() string {
	return msg.Author + " says " + msg.Body
}
