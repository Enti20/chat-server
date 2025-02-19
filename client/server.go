package client

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var (
	clientIDCounter int
	idMutex         sync.Mutex
)

// generateClientID генерирует уникальный ID для клиента.
func generateClientID() int {
	idMutex.Lock()
	defer idMutex.Unlock()
	clientIDCounter++
	return clientIDCounter
}

// Server представляет WebSocket сервер.
type Server struct {
	Pattern   string
	Messages  []*Message
	Clients   map[int]*Client
	AddCh     chan *Client
	DelCh     chan *Client
	SendAllCh chan *Message
	DoneCh    chan bool
	ErrCh     chan error
}

// Listen запускает WebSocket сервер.
func (s *Server) Listen() {
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.ErrCh <- err
			}
		}()

		client := NewClient(ws, s)
		s.Add(client)
		client.Listen()
	}

	http.HandleFunc(s.Pattern, func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		}
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Ошибка обновления соединения:", err)
			return
		}
		onConnected(ws)
	})

	for {
		select {
		case c := <-s.AddCh:
			s.Clients[c.ID] = c
			s.sendPastMessages(c)
		case c := <-s.DelCh:
			delete(s.Clients, c.ID)
			log.Printf("Клиент %d отключен", c.ID)
		case msg := <-s.SendAllCh:
			s.Messages = append(s.Messages, msg)
			s.SendAll(msg)
		case err := <-s.ErrCh:
			log.Println("Error:", err.Error())
		case <-s.DoneCh:
			return
		}
	}
}

// NewClient создает нового WebSocket клиента.
func NewClient(ws *websocket.Conn, server *Server) *Client {
	return &Client{
		ID:     generateClientID(),
		WS:     ws,
		Server: server,
		Ch:     make(chan *Message),
		DoneCh: make(chan bool),
	}
}

// SendAll отправляет сообщение всем подключенным клиентам.
func (s *Server) SendAll(msg *Message) {
	for _, client := range s.Clients {
		client.Ch <- msg // Отправляем сообщение в канал клиента
	}
}

// sendPastMessages отправляет прошлые сообщения новому клиенту.
func (s *Server) sendPastMessages(c *Client) {
	for _, msg := range s.Messages {
		c.Ch <- msg
	}
}

// Err обрабатывает ошибки.
func (s *Server) Err(err error) {
	s.ErrCh <- err
}

// Del удаляет клиента из сервера.
func (s *Server) Del(c *Client) {
	s.DelCh <- c
}

// Add добавляет нового клиента на сервер.
func (s *Server) Add(c *Client) {
	s.AddCh <- c
}
