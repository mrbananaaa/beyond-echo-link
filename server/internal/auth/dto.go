package auth

import (
	"time"

	"github.com/google/uuid"
)

type RegisterUserInput struct {
	Email          string
	Username       string
	Password       string
	Bio            string
	ProfilePicture string
}

type RegisterUserOutput struct {
	ID             uuid.UUID
	Email          string
	Username       string
	LookupID       string
	Bio            string
	ProfilePicture string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
