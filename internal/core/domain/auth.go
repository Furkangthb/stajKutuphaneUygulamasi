package domain

import (
	"context"
	"time"
)

type AuthRepository interface {
	BlackListToken(ctx context.Context, token string, duration time.Duration) error
	IsTokenBlackList(ctx context.Context, token string) (bool, error)
}
