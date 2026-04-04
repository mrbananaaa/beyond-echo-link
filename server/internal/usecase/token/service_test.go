package token

import (
	"testing"
	"time"

	redisinfra "github.com/mrbananaaa/bel-server/internal/infra/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenService_GenerateRefreshToken(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "success return refresh token",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// FIX: change this to mock
			client := redis.NewClient(&redis.Options{
				Addr: "localhost:6379",
			})
			rdb := redisinfra.New(client, "test", 2*time.Second)
			tokenStore := redisinfra.NewTokenStore(rdb)

			service := NewTokenService(
				tokenStore,
				"dummyjwtsecret",
				"bel-backend-test",
				15*time.Minute,
			)

			refreshToken, err := service.GenerateRefreshToken(t.Context(), "asdasd")
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, refreshToken)
			} else {
				require.NoError(t, err)
				require.NotNil(t, refreshToken)

				t.Logf("Refresh token: %s", refreshToken)
			}
		})
	}
}
