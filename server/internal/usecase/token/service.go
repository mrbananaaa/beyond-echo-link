package token

import (
	"time"

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

// TODO: Generate JWT

// TODO: Validate JWT

// TODO: Generate Refresh Token
