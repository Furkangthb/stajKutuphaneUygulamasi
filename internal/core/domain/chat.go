package domain

import (
	"context"
	"time"
)

type ChatMessage struct {
	ID        int
	UserID    int
	Role      string
	Message   string
	CreatedAt time.Time
}

type ChatRepository interface {
	Save(ctx context.Context, msg ChatMessage) error
	GetHistory(ctx context.Context, userID int, limit int) ([]ChatMessage, error)
}
