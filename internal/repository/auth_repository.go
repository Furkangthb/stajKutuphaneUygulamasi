package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthRepository struct {
	client *redis.Client
}

func NewAuthRepository(client *redis.Client) *AuthRepository {
	return &AuthRepository{client: client}
}

func (r *AuthRepository) BlackListToken(ctx context.Context, signature string, expiration time.Duration) error {
	return r.client.Set(ctx, "Blacklist"+signature, "true", expiration).Err()
}

func (r *AuthRepository) IsTokenBlackList(ctx context.Context, signature string) (bool, error) {
	count, err := r.client.Exists(ctx, "Blacklist"+signature).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
