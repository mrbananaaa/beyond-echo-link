package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mrbananaaa/bel-server/internal/db"
	queries "github.com/mrbananaaa/bel-server/internal/db/sqlc"
	"github.com/mrbananaaa/bel-server/internal/user"
)

type AuthService struct {
	txManager *db.TxManager
	userRepo  *user.UserRepository
}

func NewAuthService(
	txManager *db.TxManager,
	userRepo *user.UserRepository,
) *AuthService {
	return &AuthService{
		txManager: txManager,
		userRepo:  userRepo,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, error) {
	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("couldn't hash password: %w", err)
	}

	lookupID, err := generateLookupID()
	if err != nil {
		return nil, fmt.Errorf("couldn't generate lookupID: %w", err)
	}

	user, err := s.userRepo.CreateUser(ctx, queries.CreateUserParams{
		ID:       uuid.New(),
		Email:    input.Email,
		Username: input.Username,
		Password: passwordHash,
		LookupID: lookupID,
		Bio: pgtype.Text{
			String: input.Bio,
		},
		ProfilePicture: pgtype.Text{
			String: input.ProfilePicture,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &RegisterUserOutput{
		ID:             user.ID,
		Email:          user.Email,
		Username:       user.Username,
		LookupID:       user.LookupID,
		Bio:            user.Bio.String,
		ProfilePicture: user.ProfilePicture.String,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}, nil
}
