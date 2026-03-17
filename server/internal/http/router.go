package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mrbananaaa/bel-server/internal/http/handlers"
)

type Handlers struct {
	Health *handlers.HealthHandler
}

func NewRouter(h Handlers) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/health", h.Health.Routes())

	return r
}
