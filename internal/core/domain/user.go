package domain

import (
	"context"
	"time"
)

type User struct {
	ID           int       `db:"id" json:"id"`
	FirstName    string    `db:"first_name" json:"first_name"`
	LastName     string    `db:"last_name" json:"last_name"`
	Phone        string    `db:"phone" json:"phone"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Role         string    `db:"role" json:"role"`
	Created_at   time.Time `db:"created_at" json:"created_at"`
}

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context, limit, offset int) ([]*User, error)
	UpdatePassword(ctx context.Context, id int64, passwordHash string) error
}
