package httpx

import (
	"context"
	"errors"

	"github.com/mrbananaaa/bel-server/internal/domain/apperror"
)

type UserIDKey struct{}

func GetUserIDFromCtx(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey{}).(string)
	return userID, ok
}

func SetUserIDCtx(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey{}, userID)
}

type LogUserIDKey struct{}

type LogUserIDCtx struct {
	UserID string
}

func LogUserIDFromCtx(ctx context.Context) string {
	if userIDContext, ok := ctx.Value(LogUserIDKey{}).(*LogUserIDCtx); ok {
		return userIDContext.UserID
	}

	return ""
}

func SetLogUserIDCtx(ctx context.Context, userID string) error {
	if userIDContext, ok := ctx.Value(LogUserIDKey{}).(*LogUserIDCtx); ok {
		userIDContext.UserID = userID
		return nil
	}

	return apperror.InvalidCredentials(
		apperror.TypeInfrastructure,
		"couldn't set userID context value",
		errors.New("failed to set userID context value"),
	)
}
