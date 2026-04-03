package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mrbananaaa/bel-server/internal/infra/http/handlers"
	auth "github.com/mrbananaaa/bel-server/internal/infra/http/handlers/auth"
	"github.com/mrbananaaa/bel-server/internal/infra/http/middlewares"
	"github.com/mrbananaaa/bel-server/internal/infra/http/response"
)

type Handlers struct {
	Health *handlers.HealthHandler
	Auth   *auth.AuthHandler
}

type Middlewares struct {
	Log  *middlewares.LogMiddleware
	Auth *middlewares.AuthMiddleware
}

func NewRouter(h Handlers, m Middlewares) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(m.Log.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/auth", h.Auth.Routes())

	r.Mount("/health", h.Health.Routes())

	// TEST: auth middleware test endpoint
	r.Group(func(u chi.Router) {
		u.Use(m.Auth.VerifyAccessToken)
		u.Get("/protected-route", func(w http.ResponseWriter, r *http.Request) {
			response.OK(w, map[string]any{
				"message": "protected route",
			})
		})
	})

	return r
}
