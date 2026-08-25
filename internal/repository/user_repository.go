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

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (first_name,last_name,phone,email,password_hash,role,created_at) VALUES ($1,$2,$3,$4,$5,$6,now()) RETURNING id, created_at`
	err := r.db.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Phone, user.Email, user.PasswordHash, user.Role).Scan(&user.ID, &user.Created_at)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, first_name, last_name, phone, email, password_hash, role, created_at
              FROM "users"
              WHERE email = $1`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Phone, &user.Email, &user.PasswordHash, &user.Role, &user.Created_at,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	// Not: eski sorguda phone/email seçiliyordu ama Scan'e hiç verilmemişti (kolon sayısı uyuşmazlığı,
	// runtime'da "sql: expected N destination arguments in Scan" hatası verirdi) ve created_at hiç seçilmiyordu.
	query := `SELECT id,first_name,last_name,phone,email,password_hash,role,created_at
			FROM "users"
			WHERE id=$1`

	user := &domain.User{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Phone, &user.Email, &user.PasswordHash, &user.Role, &user.Created_at,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, u *domain.User) error {
	query := `UPDATE users 
			SET first_name=$1,last_name=$2,email=$3,role=$4
			WHERE id=$5
			`
	result, err := r.db.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.Role, u.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE 
			FROM users
			WHERE id=$1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context, limit, ofset int) ([]*domain.User, error) {

	query := `SELECT id,first_name,last_name,phone,email,password_hash,role,created_at
			FROM users
			ORDER BY id
			LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, ofset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User

	for rows.Next() {
		u := &domain.User{}
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Phone, &u.Email, &u.PasswordHash, &u.Role, &u.Created_at); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
