package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mrbananaaa/bel-server/internal/apperror"
	"github.com/mrbananaaa/bel-server/internal/db"
	queries "github.com/mrbananaaa/bel-server/internal/db/sqlc"
	"github.com/mrbananaaa/bel-server/internal/logger"
	"github.com/mrbananaaa/bel-server/internal/user"
)

type AuthService struct {
	txManager *db.TxManager
	userRepo  *user.UserRepository
	log       *slog.Logger
}

func NewAuthService(
	txManager *db.TxManager,
	userRepo *user.UserRepository,
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

// TODO: Login service
func (s *AuthService) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	l := s.getLogger(ctx)

	user, err := s.userRepo.GetUserByUsername(ctx, input.Username)
	if err != nil {
		logger.ErrorEvent(l,
			"user_login_failed",
			"failed to get user",
			err,
		)
		return nil, apperror.InvalidCredentials("invalid username or password", err)
	}

	if !CompareHash(input.Password, user.Password) {
		logger.ErrorEvent(l,
			"user_login_failed",
			"failed to compare hash password",
			nil,
		)
		return nil, apperror.InvalidCredentials("invalid username or password", err)
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

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// TODO: refactor this later using env variable
var (
	secretKey = []byte("thisisthejwtsupersecret")
	iss       = "bel-backend-dev"
	ttl       = 15 * time.Minute
)

func (s *AuthService) GenerateAccessToken(userID uuid.UUID) (string, error) {
	now := time.Now()
	uid := userID.String()

	claims := Claims{
		UserID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    iss,
			Subject:   uid,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func (s *AuthService) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperror.InvalidCredentials("invalid token", errors.New("failed to parse token"))
		}
		return secretKey, nil
	})
	if err != nil {
		return "", apperror.InvalidCredentials("invalid token", fmt.Errorf("failed to parse token: %w", err))
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", apperror.InvalidCredentials("invalid token", err)
	}

	return claims.UserID, nil
}

func (s *AuthService) getLogger(ctx context.Context) *slog.Logger {
	l := logger.FromContext(ctx).With("domain", "auth")
	if l == nil {
		l = s.log
	}
	return l
}
