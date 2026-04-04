package token

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mrbananaaa/bel-server/internal/apperror"
	"github.com/mrbananaaa/bel-server/internal/infra/redis"
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

// GenerateRefreshToken generate random string with rand.Read
func (s *TokenService) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", apperror.Internal(
			apperror.TypeBusiness,
			fmt.Errorf("failed to generate random byte: %w", err),
		)
	}

	return hex.EncodeToString(b), nil
}
