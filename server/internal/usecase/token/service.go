package token

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mrbananaaa/bel-server/internal/apperror"
	"github.com/mrbananaaa/bel-server/internal/infra/redis"
	"github.com/mrbananaaa/bel-server/internal/logger"
)

type TokenService struct {
	tokenStore redis.TokenStore
	jwtSecret  []byte
	jwtIss     string
	jwtTTL     time.Duration
}

func NewTokenService(
	tokenStore redis.TokenStore,
	jwtSecret string,
	jwtIss string,
	jwtTTL time.Duration,
) *TokenService {
	return &TokenService{
		tokenStore: tokenStore,
		jwtSecret:  []byte(jwtSecret),
		jwtIss:     jwtIss,
		jwtTTL:     jwtTTL,
	}
}

// JwtClaims add UserID to jwt claims
type JwtClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateJWT generate jwt with userID -
// returning tokenStr and error
func (s *TokenService) GenerateJWT(userID string) (string, error) {
	now := time.Now()

	claims := JwtClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.jwtIss,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", apperror.Internal(
			apperror.TypeBusiness,
			fmt.Errorf("failed to sign jwt: %w", err),
		)
	}

	return tokenStr, nil
}

// ValidateJWT validate tokenStr -
// returning userID and error
func (s *TokenService) ValidateJWT(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperror.InvalidCredentials(
				apperror.TypeBusiness,
				"invalid token",
				errors.New("failed to parse token: invalid token/method"),
			)
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return "", apperror.InvalidCredentials(
			apperror.TypeBusiness,
			"invalid token",
			fmt.Errorf("failed to parse token: %w", err),
		)
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok || !token.Valid {
		return "", apperror.InvalidCredentials(
			apperror.TypeBusiness,
			"invalid/expired token",
			err,
		)
	}

	return claims.UserID, nil
}

// GenerateRefreshToken generate random string with rand.Read and store it on redis
func (s *TokenService) GenerateRefreshToken(ctx context.Context, userID string) (string, error) {
	l := getLogger(ctx)

	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		err = apperror.Internal(
			apperror.TypeBusiness,
			fmt.Errorf("failed to generate random byte: %w", err),
		)
		logger.ErrorEvent(l,
			"refreshtoken_generation_failed",
			"failed to generate token from rand.Read",
			err,
		)

		return "", err
	}
	refreshToken := hex.EncodeToString(b)

	if err = s.tokenStore.SetRefreshToken(ctx, refreshToken, userID, 7*24*time.Hour); err != nil {
		err := apperror.Internal(
			apperror.TypeBusiness,
			fmt.Errorf("failed to set refresh token: %w", err),
		)
		logger.ErrorEvent(l,
			"refreshtoken_generation_failed",
			"failed to save refresh token to redis",
			err,
		)

		return "", err
	}

	return refreshToken, nil
}

func (s *TokenService) GetUserIDFromRefreshToken(ctx context.Context, token string) (string, error) {
	l := getLogger(ctx)

	userID, exists, err := s.tokenStore.GetUserIDByToken(ctx, token)
	if err != nil {
		err := apperror.Internal(
			apperror.TypeBusiness,
			fmt.Errorf("failed to get userID with token: %w", err),
		)
		logger.ErrorEvent(l,
			"refreshtoken_find_failed",
			"failed to get userID with token",
			err,
		)

		return "", err
	}

	if !exists {
		err := apperror.Internal(
			apperror.TypeBusiness,
			fmt.Errorf("token not exists"),
		)
		logger.ErrorEvent(l,
			"refreshtoken_find_failed",
			"token not exists",
			err,
		)

		return "", err
	}

	return userID, nil
}

func (s *TokenService) RotateRefreshToken(
	ctx context.Context,
	oldToken string,
	newToken string,
	userID string,
) error {
	l := getLogger(ctx)

	if err := s.tokenStore.RotateRefreshToken(ctx, oldToken, newToken, userID, 7*24*time.Hour); err != nil {
		err := apperror.Internal(
			apperror.TypeBusiness,
			fmt.Errorf("failed to rotate refresh token to redis: %w", err),
		)
		logger.ErrorEvent(l,
			"refreshtoken_rotation_failed",
			"failed to rotate refresh token",
			err,
		)

		return err
	}

	return nil
}

func (s *TokenService) DeleteRefreshToken(ctx context.Context, token string) error {
	l := logger.FromContext(ctx)

	if err := s.tokenStore.DeleteRefreshToken(ctx, token); err != nil {
		err := apperror.Internal(
			apperror.TypeBusiness,
			fmt.Errorf("failed to delete refresh token: %w", err),
		)
		logger.ErrorEvent(l,
			"refreshtoken_deletion_failed",
			"failed to delete refresh token from redis",
			err,
		)

		return err
	}

	return nil
}

func getLogger(ctx context.Context) *slog.Logger {
	return logger.FromContext(ctx).With("domain", "token")
}
