package httpx

import (
	"context"
	"errors"

	"github.com/mrbananaaa/bel-server/internal/apperror"
)

type UserIDKey struct{}

type UserIDContext struct {
	UserID string
}

func UserIDFromContext(ctx context.Context) string {
	if userIDContext, ok := ctx.Value(UserIDKey{}).(*UserIDContext); ok {
		return userIDContext.UserID
	}

	return ""
}

func SetUserIDContext(ctx context.Context, userID string) error {
	if userIDContext, ok := ctx.Value(UserIDKey{}).(*UserIDContext); ok {
		userIDContext.UserID = userID
		return nil
	}

	return apperror.InvalidCredentials(
		apperror.TypeInfrastructure,
		"couldn't set userID context value",
		errors.New("failed to set userID context value"),
	)
}
