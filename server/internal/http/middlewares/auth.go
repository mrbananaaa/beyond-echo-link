package middlewares

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/mrbananaaa/bel-server/internal/apperror"
	"github.com/mrbananaaa/bel-server/internal/auth"
	"github.com/mrbananaaa/bel-server/internal/http/httpx"
	"github.com/mrbananaaa/bel-server/internal/http/response"
	"github.com/mrbananaaa/bel-server/internal/logger"
)

type AuthMiddleware struct {
	authService *auth.AuthService
	log         *slog.Logger
}

func NewAuthMiddleware(
	authService *auth.AuthService,
	log *slog.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		log:         log,
	}
}

func (m *AuthMiddleware) VerifyAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			err := apperror.InvalidCredentials(apperror.TypeInfrastructure, "invalid token", errors.New("empty Authorization header"))
			logger.ErrorEvent(m.log,
				"token_validation_failed",
				"empty Authorization header",
				err,
			)
			response.Error(w, r, err)
			return
		}

		authHeaderValues := strings.Split(authHeader, " ")
		if len(authHeaderValues) != 2 || strings.ToLower(authHeaderValues[0]) != "bearer" {
			err := apperror.InvalidCredentials(apperror.TypeInfrastructure, "invalid token", errors.New("invalid Authorization header"))
			logger.ErrorEvent(m.log,
				"token_validation_failed",
				"invalid Authorization header",
				err,
			)
			response.Error(w, r, err)
			return
		}

		userID, err := m.authService.ValidateToken(authHeaderValues[1])
		if err != nil {
			logger.ErrorEvent(m.log,
				"token_validation_failed",
				err.Error(),
				err,
			)
			response.Error(w, r, err)
			return
		}

		if err := httpx.SetUserIDContext(r.Context(), userID); err != nil {
			// FIX: change to app error
			logger.ErrorEvent(m.log,
				"token_validation_failed",
				"couldn't set userID context value",
				err,
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func MockAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIDHeader := r.Header.Get("X-User-ID")

			ctx := r.Context()

			if userIDHeader != "" {
				ctx = logger.WithUserID(ctx, userIDHeader)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
