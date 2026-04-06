package logger

import (
	"context"
	"log/slog"

	"github.com/mrbananaaa/bel-server/internal/infra/http/httpx"
)

type ctxKey struct{}

func WithContext(
	ctx context.Context,
	logger *slog.Logger,
) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return l
	}

	return slog.Default()
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, httpx.LogUserIDKey{}, userID)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(httpx.LogUserIDKey{}).(string)
	return id, ok
}
