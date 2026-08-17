package domain

import "time"

type Reservation struct {
	ID         int       `db:"id" json:"id"`
	UserID     int       `db:"user_id" json:"user_id" `
	BookID     int       `db:"book_id" json:"book_id"`
	Status     string    `db:"status" json:"status"`
	ReservedAt time.Time `db:"reserved_at" json:"reserved_at"`
	DueDate    time.Time `db:"due_date" json:"due_date"`
}
