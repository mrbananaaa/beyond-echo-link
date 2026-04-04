package redis

import (
	"context"
	"time"
)

type TokenStore interface {
	SetRefreshToken(ctx context.Context, token string, userID string, ttl time.Duration) error
	GetUserIDByToken(ctx context.Context, token string) (string, bool, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	RotateRefreshToken(ctx context.Context, oldToken string, newToken string, userID string, ttl time.Duration) error
}

type RedisTokenStore struct {
	redis *Redis
}

func NewTokenStore(r *Redis) *RedisTokenStore {
	return &RedisTokenStore{
		redis: r,
	}
}

func (s *RedisTokenStore) key(token string) string {
	return s.redis.Key("auth", "refresh", token)
}

func (s *RedisTokenStore) SetRefreshToken(
	ctx context.Context,
	token string,
	userID string,
	ttl time.Duration,
) error {
	key := s.key(token)

	return s.redis.Set(ctx, key, userID, ttl)
}

func (s *RedisTokenStore) GetUserIDByToken(
	ctx context.Context,
	token string,
) (string, bool, error) {
	key := s.key(token)

	val, err := s.redis.Get(ctx, key)
	if err != nil {
		return "", false, err
	}

	if val == "" {
		return "", false, nil
	}

	return val, true, nil
}

func (s *RedisTokenStore) DeleteRefreshToken(
	ctx context.Context,
	token string,
) error {
	key := s.key(token)

	return s.redis.Delete(ctx, key)
}

func (s *RedisTokenStore) RotateRefreshToken(
	ctx context.Context,
	oldToken string,
	newToken string,
	userID string,
	ttl time.Duration,
) error {
	if err := s.DeleteRefreshToken(ctx, oldToken); err != nil {
		return err
	}

	return s.SetRefreshToken(ctx, newToken, userID, ttl)
}
