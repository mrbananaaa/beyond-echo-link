package websocket

import (
	"context"
	"fmt"

	ws "github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/mrbananaaa/bel-server/internal/domain/chat"
	"github.com/mrbananaaa/bel-server/internal/logger"
)

type Client struct {
	ID   string
	Conn *ws.Conn
	Send chan []byte

	Rooms map[string]bool
}

func (c *Client) ReadPump(ctx context.Context, hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close(ws.StatusNormalClosure, "closed")
	}()

	l := logger.FromContext(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			var msg chat.Message
			err := wsjson.Read(ctx, c.Conn, &msg)
			if err != nil {
				logger.ErrorEvent(l,
					"msg_decode_failed",
					"err decoding websocket message",
					fmt.Errorf("couldn't decode websocket message: %w", err),
				)
				return
			}

			switch msg.Type {
			case chat.JoinRoom:
				hub.JoinRoom <- JoinRequest{
					Client: c,
					RoomID: msg.RoomID,
				}
			case chat.LeaveRoom:
				hub.LeaveRoom <- LeaveRequest{
					Client: c,
					RoomID: msg.RoomID,
				}
			case chat.MessageText:
				hub.Broadcast <- BroadcastRequest{
					RoomID: msg.RoomID,
					Data:   msg.Payload,
				}
			}
		}
	}
}

func (c *Client) WritePump(ctx context.Context) {
	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				return
			}

			err := wsjson.Write(ctx, c.Conn, msg)
			if err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
