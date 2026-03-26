package auth

import "github.com/go-chi/chi/v5"

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/signup", h.SignUp)

	return r
}
