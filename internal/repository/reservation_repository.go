package repository

import (
	"context"
	"database/sql"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
)

type ReservationRepository struct {
	db *sql.DB
}

func NewResservationRepository(db *sql.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (r *ReservationRepository) ReservationCreate(ctx context.Context, reservation *domain.Reservation) error {
	query := `INSERT INTO reservations( user_id,book_id,status,due_date) VALUES($1,$2,$3,$4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, reservation.UserID, reservation.BookID, reservation.Status, reservation.DueDate).Scan(&reservation.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReservationRepository) ReservationUpdate(ctx context.Context, reservation *domain.Reservation) error {
	query := `UPDATE reservations
			SET status=$1
			WHERE id=$2`
	result, err := r.db.ExecContext(ctx, query, reservation.Status, reservation.ID)
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

func (r *ReservationRepository) ReservationGetByID(ctx context.Context, id int) (*domain.Reservation, error) {
	query := `SELECT id, user_id, book_id, status, reserved_at, due_date
			FROM reservations
			WHERE id=$1`

	reservation := &domain.Reservation{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&reservation.ID, &reservation.UserID, &reservation.BookID,
		&reservation.Status, &reservation.ReservedAt, &reservation.DueDate,
	)
	if err != nil {
		return nil, err
	}
	return reservation, nil
}

func (r *ReservationRepository) ReservationGetUserByID(ctx context.Context, userID int) ([]*domain.ReservationWithUser, error) {
	query := `SELECT r.id, r.user_id, r.book_id, r.status, r.reserved_at, r.due_date, u.first_name, u.last_name
			FROM reservations AS r
			INNER JOIN users AS u ON r.user_id=u.id
			WHERE u.id=$1`

	reservations := []*domain.ReservationWithUser{}
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		reserve := &domain.ReservationWithUser{}
		err := rows.Scan(
			&reserve.ID, &reserve.UserID, &reserve.BookID, &reserve.Status,
			&reserve.ReservedAt, &reserve.DueDate, &reserve.FirstName, &reserve.LastName,
		)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reserve)
	}
	return reservations, nil
}
