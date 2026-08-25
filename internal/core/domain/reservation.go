package domain

import (
	"context"
	"time"
)

type Reservation struct {
	ID         int       `db:"id" json:"id"`
	UserID     int       `db:"user_id" json:"user_id" `
	BookID     int       `db:"book_id" json:"book_id"`
	Status     string    `db:"status" json:"status"`
	ReservedAt time.Time `db:"reserved_at" json:"reserved_at"`
	DueDate    time.Time `db:"due_date" json:"due_date"`
}

type ReservationWithUser struct {
	ID         int       `db:"id" json:"id"`
	UserID     int       `db:"user_id" json:"user_id"`
	BookID     int       `db:"book_id" json:"book_id"`
	Status     string    `db:"status" json:"status"`
	ReservedAt time.Time `db:"reserved_at" json:"reserved_at"`
	DueDate    time.Time `db:"due_date" json:"due_date"`
	FirstName  string    `db:"first_name" json:"first_name"`
	LastName   string    `db:"last_name" json:"last_name"`
}

type ReservationFull struct {
	ID         int       `db:"id" json:"id"`
	UserID     int       `db:"user_id" json:"user_id"`
	BookID     int       `db:"book_id" json:"book_id"`
	Status     string    `db:"status" json:"status"`
	ReservedAt time.Time `db:"reserved_at" json:"reserved_at"`
	DueDate    time.Time `db:"due_date" json:"due_date"`
	FirstName  string    `db:"first_name" json:"first_name"`
	LastName   string    `db:"last_name" json:"last_name"`
	BookTitle  string    `db:"book_title" json:"book_title"`
	BookAuthor string    `db:"book_author" json:"book_author"`
}

type ReservationRepository interface {
	Create(ctx context.Context, r *Reservation) error
	Update(ctx context.Context, id int, status string) error
	GetByUserID(ctx context.Context, userId int) ([]*ReservationWithUser, error)
}
