package user

import (
	"time"

	"github.com/google/uuid"
	queries "github.com/mrbananaaa/bel-server/internal/infra/db/sqlc"
)

type User struct {
	ID             uuid.UUID
	Email          string
	Username       string
	Password       string
	LookupID       string
	Bio            string
	ProfilePicture string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (u *User) FromDB(user queries.User) {
	u.ID = user.ID
	u.Email = user.Email
	u.Username = user.Username
	u.Password = user.Password
	u.LookupID = user.LookupID
	u.Bio = user.Bio.String
	u.ProfilePicture = user.ProfilePicture.String
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
}
