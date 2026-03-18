package logger

import (
	"context"
	"log/slog"
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

type userKey struct{}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userKey{}, userID)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userKey{}).(string)
	return id, ok
}
