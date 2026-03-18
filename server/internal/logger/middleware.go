package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func Middleware(base *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			reqID := uuid.NewString()

			rw := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			l := base.With(
				"request_id", reqID,
				"method", r.Method,
				"path", r.URL.Path,
			)

			userID, ok := UserIDFromContext(r.Context())
			if !ok {
				userID = ""
			}
			l = l.With("user_id", userID)

			ctx := WithContext(r.Context(), l)

			next.ServeHTTP(rw, r.WithContext(ctx))

			duration := time.Since(start).Milliseconds()

			l.Info("request completed",
				"event", "http_request_completed",
				"status", rw.status,
				"duration_ms", duration,
			)
		})
	}
}
