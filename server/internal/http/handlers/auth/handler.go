package auth

import (
	"encoding/json"
	"net/http"

	"github.com/mrbananaaa/bel-server/internal/apperror"
	"github.com/mrbananaaa/bel-server/internal/auth"
	"github.com/mrbananaaa/bel-server/internal/http/response"
	"github.com/mrbananaaa/bel-server/internal/logger"
	"github.com/mrbananaaa/bel-server/internal/validation"
)

type AuthHandler struct {
	validator   *validation.Validator
	authService *auth.AuthService
}

func NewAuthHandler(
	validator *validation.Validator,
	authService *auth.AuthService,
) *AuthHandler {
	return &AuthHandler{
		validator:   validator,
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

	if field, err := h.validator.Validate(req); err != nil {
		logger.ErrorValidation(l, err)
		response.Error(w, r, apperror.ValidationError(field...))
		return
	}

	// TODO: register
	// TODO: token
	// TODO: response
	response.JSON(w, http.StatusOK, map[string]any{"message": "sign up route"})
}
