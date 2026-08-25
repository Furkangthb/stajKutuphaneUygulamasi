package repository

import (
	"context"
	"database/sql"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
)

type chatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *chatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) Save(ctx context.Context, msg domain.ChatMessage) error {
	query := `INSERT INTO chat_history (user_id, role, message) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, msg.UserID, msg.Role, msg.Message)
	return err
}

func (r *chatRepository) GetHistory(ctx context.Context, userID int, limit int) ([]domain.ChatMessage, error) {
	query := `SELECT id, user_id, role, message, created_at FROM chat_history WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2`
	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sonuc []domain.ChatMessage
	for rows.Next() {
		var m domain.ChatMessage
		if err := rows.Scan(&m.ID, &m.UserID, &m.Role, &m.Message, &m.CreatedAt); err != nil {
			return nil, err
		}
		sonuc = append(sonuc, m)
	}
	return sonuc, nil
}
