package middlewares

import (
	"net/http"

	"github.com/mrbananaaa/bel-server/internal/logger"
)

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
