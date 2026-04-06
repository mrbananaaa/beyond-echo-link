package ws

import (
	"net/http"

	"github.com/mrbananaaa/bel-server/internal/infra/http/response"
	"github.com/mrbananaaa/bel-server/internal/logger"
)

type WsHandler struct {
}

func NewWsHandler() *WsHandler {
	return &WsHandler{}
}

// Websocket /ws
func (h *WsHandler) Websocket(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())

	logger.InfoEvent(l, "websocket_connect", "connecting to websocket")

	response.OK(w, map[string]any{
		"message": "ws route",
	})
}
