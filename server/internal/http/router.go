package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mrbananaaa/bel-server/internal/http/handlers"
	auth "github.com/mrbananaaa/bel-server/internal/http/handlers/auth"
	"github.com/mrbananaaa/bel-server/internal/http/middlewares"
	"github.com/mrbananaaa/bel-server/internal/logger"
)

type Handlers struct {
	Health *handlers.HealthHandler
	Auth   *auth.AuthHandler
}

func NewRouter(h Handlers) *chi.Mux {
	log := logger.New(logger.Config{
		Env:     "dev",
		Service: "api",
	})

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middlewares.MockAuth())
	r.Use(logger.Middleware(log))
	r.Use(middleware.Recoverer)

	r.Mount("/auth", h.Auth.Routes())

	r.Mount("/health", h.Health.Routes())

	return r
}
