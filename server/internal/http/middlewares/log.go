package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/mrbananaaa/bel-server/internal/http/httpx"
	"github.com/mrbananaaa/bel-server/internal/logger"
)

type LogMiddleware struct {
	log *slog.Logger
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func NewLogMiddleware(
	log *slog.Logger,
) *LogMiddleware {
	return &LogMiddleware{
		log: log,
	}
}

func (m *LogMiddleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// RequestID from chi middleware
		var reqID string
		if reqIDFromCtx, ok := r.Context().Value(middleware.RequestIDKey).(string); ok {
			reqID = reqIDFromCtx
		}

		// UserID for logging
		userIDCtx := &httpx.UserIDContext{UserID: ""}

		// keep track status
		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		l := m.log.With(
			"request_id", reqID,
			"method", r.Method,
			"path", r.URL.Path,
		)

		ctx := context.WithValue(r.Context(), httpx.UserIDKey{}, userIDCtx)
		ctx = logger.WithContext(ctx, l)

		next.ServeHTTP(rw, r.WithContext(ctx))

		duration := time.Since(start).Milliseconds()

		if userIDCtx.UserID != "" {
			l = l.With("user_id", userIDCtx.UserID)
		}

		l.Info("request completed",
			"event", "http_request_completed",
			"status", rw.status,
			"duration_ms", duration,
		)
	})
}
