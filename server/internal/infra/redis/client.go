package redis

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client  *redis.Client
	prefix  string
	timeout time.Duration
}

func New(
	client *redis.Client,
	prefix string,
	timeout time.Duration,
) *Redis {
	return &Redis{
		client:  client,
		prefix:  prefix,
		timeout: timeout,
	}
}

func (r *Redis) Close() error {
	return r.client.Close()
}

func (r *Redis) Ping(ctx context.Context) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	return r.client.Ping(ctx).Err()
}

func (r *Redis) buildKey(parts ...string) string {
	key := strings.Join(parts, ":")

	if r.prefix != "" {
		return r.prefix + ":" + key
	}

	return key
}

func (r *Redis) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, r.timeout)
}

func (r *Redis) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *Redis) GetWithExists(ctx context.Context, key string) (string, bool, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	return val, true, err
}

func (r *Redis) Delete(ctx context.Context, keys ...string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	return r.client.Del(ctx, keys...).Err()
}

func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	n, err := r.client.Exists(ctx, key).Result()
	return n > 0, err
}

func (r *Redis) Expire(ctx context.Context, key string, ttl time.Duration) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	return r.client.Expire(ctx, key, ttl).Err()
}

func (r *Redis) Key(parts ...string) string {
	return r.buildKey(parts...)
}
