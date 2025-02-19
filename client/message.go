package client

// Message представляет сообщение чата.
type Message struct {
	ClientID int    `json:"client_id"`
	Author   string `json:"author"`
	Body     string `json:"body"`
}
