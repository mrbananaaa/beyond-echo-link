package ws

import (
	"errors"
	"fmt"
	"net/http"

	ws "github.com/coder/websocket"
	"github.com/mrbananaaa/bel-server/internal/domain/apperror"
	"github.com/mrbananaaa/bel-server/internal/domain/token"
	"github.com/mrbananaaa/bel-server/internal/infra/http/response"
	"github.com/mrbananaaa/bel-server/internal/infra/websocket"
	"github.com/mrbananaaa/bel-server/internal/logger"
)

type WsHandler struct {
	tokenService *token.TokenService
}

func NewWsHandler(
	tokenService *token.TokenService,
) *WsHandler {
	return &WsHandler{
		tokenService: tokenService,
	}
}

// Websocket /ws
func (h *WsHandler) Websocket(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())

	token := r.URL.Query().Get("token")
	if token == "" {
		err := apperror.InvalidCredentials(
			apperror.TypeInfrastructure,
			"invalid token",
			errors.New("couldn't get token from query"),
		)
		logger.ErrorEvent(l,
			"ws_upgrade_failed",
			"couldn't get token from query",
			err,
		)
		response.Error(w, r, err)
		return
	}
	userID, err := h.tokenService.ValidateJWT(token)
	if err != nil {
		response.Error(w, r, err)
		return
	}

	conn, err := ws.Accept(w, r, nil)
	if err != nil {
		apperr := apperror.InvalidCredentials(
			apperror.TypeInfrastructure,
			fmt.Sprintf("socket upgrade failed: %v", err),
			fmt.Errorf("can't accept websocket connections: %w", err),
		)
		logger.ErrorEvent(l,
			"ws_upgrade_failed",
			"can't accept websocket connections",
			apperr,
		)
		response.Error(w, r, apperr)
		return
	}

	client := &websocket.Client{
		ID:    userID,
		Conn:  conn,
		Rooms: make(map[string]bool),
		Send:  make(chan []byte),
	}

	logger.InfoEvent(l,
		"client_created",
		"success to creating client",
		"client", fmt.Sprintf("%+v", client),
	)
}
