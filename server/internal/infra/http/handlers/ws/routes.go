package ws

import "github.com/go-chi/chi/v5"

func (h *WsHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.Websocket)

	return r
}
