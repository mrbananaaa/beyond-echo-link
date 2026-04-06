package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mrbananaaa/bel-server/internal/apperror"
	"github.com/mrbananaaa/bel-server/internal/domain/auth"
	"github.com/mrbananaaa/bel-server/internal/domain/token"
	"github.com/mrbananaaa/bel-server/internal/infra/http/httpx"
	"github.com/mrbananaaa/bel-server/internal/infra/http/response"
	"github.com/mrbananaaa/bel-server/internal/logger"
	"github.com/mrbananaaa/bel-server/internal/validation"
)

type AuthHandler struct {
	validator    *validation.Validator
	authService  *auth.AuthService
	tokenService *token.TokenService
}

func NewAuthHandler(
	validator *validation.Validator,
	authService *auth.AuthService,
	tokenService *token.TokenService,
) *AuthHandler {
	return &AuthHandler{
		validator:    validator,
		authService:  authService,
		tokenService: tokenService,
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

	accessToken, err := h.tokenService.GenerateJWT(user.ID.String())
	if err != nil {
		response.Error(w, r, err)
		return
	}

	refreshToken, err := h.tokenService.GenerateRefreshToken(r.Context(), user.ID.String())
	if err != nil {
		response.Error(w, r, err)
		return
	}

	httpx.SetRefreshTokenCookie(w, refreshToken)

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

	accessToken, err := h.tokenService.GenerateJWT(user.ID.String())
	if err != nil {
		response.Error(w, r, err)
		return
	}

	refreshToken, err := h.tokenService.GenerateRefreshToken(r.Context(), user.ID.String())
	if err != nil {
		response.Error(w, r, err)
		return
	}

	httpx.SetRefreshTokenCookie(w, refreshToken)

	response.JSON(w, http.StatusOK, map[string]any{
		"user_id":      user.ID,
		"access_token": accessToken,
	})
}

// SignOut /auth/signout
// FIX: use userID from context to delete token
func (h *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())

	refreshToken, err := httpx.GetRefreshTokenCookie(r)
	if err != nil {
		err := apperror.Internal(
			apperror.TypeInfrastructure,
			fmt.Errorf("failed to get refresh token from cookie: %v", err),
		)
		logger.ErrorEvent(l,
			"auth_signingout_failed",
			"failed to get refresh token from cookie",
			err,
		)

		response.Error(w, r, err)
		return
	}

	if err := h.tokenService.DeleteRefreshToken(r.Context(), refreshToken); err != nil {
		response.Error(w, r, err)
		return
	}

	httpx.ClearRefreshTokenCokie(w)

	response.OK(w, map[string]any{
		"status": "ok",
	})
}

// Refresh /auth/refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())

	refreshToken, err := httpx.GetRefreshTokenCookie(r)
	if err != nil {
		err := apperror.InvalidCredentials(
			apperror.TypeInfrastructure,
			"invalid refresh token",
			fmt.Errorf("failed to get refresh token from cookie: %w", err),
		)
		logger.ErrorEvent(l,
			"auth_refreshtoken_failed",
			"failed to get refresh token from cookie",
			err,
		)

		response.Error(w, r, err)
		return
	}

	// TODO: use userID from context instead
	userID, err := h.tokenService.GetUserIDFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		response.Error(w, r, err)
		return
	}

	accessToken, err := h.tokenService.GenerateJWT(userID)
	if err != nil {
		logger.ErrorEvent(l,
			"refreshtoken_rotation_failed",
			"failed to issuing new access token",
			err,
		)

		response.Error(w, r, err)
		return
	}

	newRefreshToken, err := h.tokenService.GenerateRefreshToken(r.Context(), userID)
	if err != nil {
		response.Error(w, r, err)
		return
	}

	if err := h.tokenService.RotateRefreshToken(r.Context(), refreshToken, newRefreshToken, userID); err != nil {
		response.Error(w, r, err)
		return
	}

	httpx.SetRefreshTokenCookie(w, newRefreshToken)

	response.OK(w, map[string]any{
		"access_token": accessToken,
	})
}
