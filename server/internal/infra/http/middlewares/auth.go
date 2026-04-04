package middlewares

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/mrbananaaa/bel-server/internal/apperror"
	"github.com/mrbananaaa/bel-server/internal/infra/http/httpx"
	"github.com/mrbananaaa/bel-server/internal/infra/http/response"
	"github.com/mrbananaaa/bel-server/internal/logger"
	"github.com/mrbananaaa/bel-server/internal/usecase/token"
)

type AuthMiddleware struct {
	tokenService *token.TokenService
	log          *slog.Logger
}

func NewAuthMiddleware(
	tokenService *token.TokenService,
	log *slog.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
		log:          log,
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

		userID, err := m.tokenService.ValidateJWT(authHeaderValues[1])
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
