package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/mrbananaaa/bel-server/internal/auth"
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(
	authService *auth.AuthService,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()

	return r
}
