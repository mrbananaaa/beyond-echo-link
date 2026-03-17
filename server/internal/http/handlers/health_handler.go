package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	resp := map[string]string{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
