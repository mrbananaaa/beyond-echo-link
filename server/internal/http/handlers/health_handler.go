package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrbananaaa/bel-server/internal/http/response"
	"github.com/mrbananaaa/bel-server/internal/logger"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.Check)

	return r
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())
	log.Info("handling request")

	resp := map[string]string{
		"status": "ok",
	}

	response.OK(w, resp)
}
