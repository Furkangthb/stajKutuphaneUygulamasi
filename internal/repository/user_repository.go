package repository

import (
	"context"
	"database/sql"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) UserRegister(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (first_name,last_name,phone,email,password_hash,role) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Phone, user.Email, user.PasswordHash, user.Role).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UserLogin(ctx context.Context, user *domain.User) error {
	query := `SELECT id,first_name,last_name,phone,email,sifre_hash
			FROM "users"
			WHERE email=$1`

	err := r.db.QueryRowContext(ctx, query, user.Email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Phone, &user.PasswordHash)
	if err != nil {
		return err

	}
	return nil
}
