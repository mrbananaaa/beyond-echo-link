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

	if err := h.validator.Validate(req); err != nil {
		logger.ErrorValidation(l, err)
		response.Error(w, r, err)
		return
	}

	// TODO: register
	user, err := h.authService.RegisterUser(r.Context(), auth.RegisterUserInput(req))
	if err != nil {
		response.Error(w, r, apperror.ErrInternal)
		return
	}

	// TODO: token
	// TODO: response
	response.JSON(w, http.StatusCreated, map[string]any{"user": user})
}
