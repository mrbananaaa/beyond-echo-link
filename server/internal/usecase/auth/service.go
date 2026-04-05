package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mrbananaaa/bel-server/internal/apperror"
	"github.com/mrbananaaa/bel-server/internal/infra/db"
	queries "github.com/mrbananaaa/bel-server/internal/infra/db/sqlc"
	"github.com/mrbananaaa/bel-server/internal/logger"
	"github.com/mrbananaaa/bel-server/internal/usecase/user"
)

type AuthService struct {
	txManager *db.TxManager
	userRepo  user.UserRepository
	log       *slog.Logger
}

func NewAuthService(
	txManager *db.TxManager,
	userRepo user.UserRepository,
	log *slog.Logger,
) *AuthService {
	return &AuthService{
		txManager: txManager,
		userRepo:  userRepo,
		log:       log.With("domain", "auth"),
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, error) {
	l := s.getLogger(ctx)

	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		return nil, apperror.Internal(apperror.TypeBusiness, fmt.Errorf("failed to hash password: %w", err))
	}

	lookupID, err := generateLookupID()
	if err != nil {
		return nil, apperror.Internal(apperror.TypeBusiness, fmt.Errorf("couldn't generate lookupID: %w", err))
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
		logger.ErrorEvent(l,
			"user_creation_failed",
			"failed to create user",
			err,
		)

		return nil, err
	}

	logger.InfoEvent(l,
		"user_created",
		"user created successfully",
		"user_id", user.ID,
	)
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

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	l := s.getLogger(ctx)

	user, err := s.userRepo.GetUserByUsername(ctx, input.Username)
	if err != nil {
		logger.ErrorEvent(l,
			"user_login_failed",
			"failed to get user",
			err,
		)
		return nil, apperror.InvalidCredentials(apperror.TypeBusiness, "invalid username or password", err)
	}

	if !CompareHash(input.Password, user.Password) {
		// INFO: apperrr to return
		err := apperror.InvalidCredentials(
			apperror.TypeBusiness,
			"invalid username or password",
			err,
		)
		logger.ErrorEvent(l,
			"user_login_failed",
			"failed to compare hash password",
			err,
		)
		return nil, err
	}

	logger.InfoEvent(l,
		"user_loggedin",
		"user login successfully",
		"user_id", user.ID,
	)
	return &LoginOutput{
		ID: user.ID,
	}, nil
}

func (s *AuthService) getLogger(ctx context.Context) *slog.Logger {
	l := logger.FromContext(ctx).With("domain", "auth")
	if l == nil {
		l = s.log
	}
	return l
}
