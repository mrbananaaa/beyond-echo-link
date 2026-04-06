package websocket

import ws "github.com/coder/websocket"

type Client struct {
	ID   string
	Conn *ws.Conn
	Send chan []byte

	Rooms map[string]bool
}
