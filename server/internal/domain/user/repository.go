package user

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Save(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*User, error)
	GetByLookupID(ctx context.Context, lookupID string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}
