package auth

import (
	"encoding/json"
	"net/http"

	"github.com/mrbananaaa/bel-server/internal/apperror"
	"github.com/mrbananaaa/bel-server/internal/infra/http/response"
	"github.com/mrbananaaa/bel-server/internal/logger"
	"github.com/mrbananaaa/bel-server/internal/usecase/auth"
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

// SignUp /auth/signup
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())

	var req signupRequest

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

	user, err := h.authService.RegisterUser(r.Context(), auth.RegisterUserInput(req))
	if err != nil {
		response.Error(w, r, err)
		return
	}

	// TODO: refresh token

	accessToken, err := h.authService.GenerateAccessToken(user.ID)
	if err != nil {
		response.Error(w, r, err)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"user":         user,
		"access_token": accessToken,
	})
}

// Signin /auth/signin
func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())

	var req signinRequest

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

	user, err := h.authService.Login(r.Context(), auth.LoginInput(req))
	if err != nil {
		response.Error(w, r, err)
		return
	}

	accessToken, err := h.authService.GenerateAccessToken(user.ID)
	if err != nil {
		response.Error(w, r, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"user_id":      user.ID,
		"access_token": accessToken,
	})
}

// RefreshToken /auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {

}
