package auth

import (
	"encoding/json"
	"net/http"

	"github.com/mrbananaaa/bel-server/internal/apperror"
	"github.com/mrbananaaa/bel-server/internal/auth"
	"github.com/mrbananaaa/bel-server/internal/http/response"
	"github.com/mrbananaaa/bel-server/internal/logger"
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

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())

	var req registerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.ErrorParseJSON(l, err)
		response.Error(w, r, apperror.ErrBadRequest)
		return
	}

	// TODO: validate

	// TODO: register
	// TODO: token
	// TODO: response
	response.JSON(w, http.StatusOK, map[string]any{"message": "sign up route"})
}
