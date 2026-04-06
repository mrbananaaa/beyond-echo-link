package middlewares

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/mrbananaaa/bel-server/internal/domain/apperror"
	"github.com/mrbananaaa/bel-server/internal/domain/token"
	"github.com/mrbananaaa/bel-server/internal/infra/http/httpx"
	"github.com/mrbananaaa/bel-server/internal/infra/http/response"
	"github.com/mrbananaaa/bel-server/internal/logger"
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

		if err := httpx.SetLogUserIDCtx(r.Context(), userID); err != nil {
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

func (m *AuthMiddleware) VerifyRefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())

		refreshToken, err := httpx.GetRefreshTokenCookie(r)
		if err != nil {
			err := apperror.InvalidCredentials(apperror.TypeInfrastructure, "invalid token", err)
			logger.ErrorEvent(l,
				"refreshtoken_validate_failed",
				"couldn't get refresh token from cookie",
				err,
			)

			response.Error(w, r, err)
			return
		}

		// TODO: USERIDCONTEXT

		logger.InfoEvent(l,
			"refresh_token_middleware",
			"refresh token middleware",
			"refresh_token", refreshToken,
		)

		next.ServeHTTP(w, r)
	})
}
