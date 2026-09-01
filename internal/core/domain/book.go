package domain

import (
	"context"
	"time"
)

type Book struct {
	ID          int       `db:"id" json:"id"`
	ISBN        string    `db:"ısbn" json:"ısbn"`
	Title       string    `db:"title" json:"title"`
	Author      string    `db:"author" json:"author"`
	Genre       string    `db:"genre" json:"genre"`
	PublishDate time.Time `db:"publish_date" json:"publish_date"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	Available   bool      `json:"available"`
}

type BookRepository interface {
	BookCreate(ctx context.Context, b *Book) error
	BookUpdate(ctx context.Context, b *Book) error
	BookDelete(ctx context.Context, id int64) error
	BookGetByID(ctx context.Context, id int64) (*Book, error)
	BookList(ctx context.Context, limit, offset int) ([]*Book, error)
	BookSearch(ctx context.Context, limit int, keywords []string) ([]*Book, error)
}
