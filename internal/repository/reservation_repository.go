package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
)

var (
	ErrBookNotFound     = errors.New("kitap bulunamadi")
	ErrBookNotAvailable = errors.New("kitap su anda musait degil")
)

type ReservationRepository struct {
	db *sql.DB
}

func NewResservationRepository(db *sql.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (r *ReservationRepository) ReservationCreate(ctx context.Context, reservation *domain.Reservation) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var bookExists bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM books WHERE id=$1 FOR UPDATE)`, reservation.BookID).Scan(&bookExists)
	if err != nil {
		return err
	}
	if !bookExists {
		return ErrBookNotFound
	}

	var alreadyReserved bool
	err = tx.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM reservations WHERE book_id=$1 AND status='active')`,
		reservation.BookID,
	).Scan(&alreadyReserved)
	if err != nil {
		return err
	}
	if alreadyReserved {
		return ErrBookNotAvailable
	}

	query := `INSERT INTO reservations( user_id,book_id,status,due_date) VALUES($1,$2,$3,$4) RETURNING id, reserved_at`
	err = tx.QueryRowContext(ctx, query, reservation.UserID, reservation.BookID, reservation.Status, reservation.DueDate).Scan(&reservation.ID, &reservation.ReservedAt)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ReservationRepository) ReservationUpdate(ctx context.Context, reservation *domain.Reservation) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var mevcutStatus string
	var bookID int
	err = tx.QueryRowContext(ctx, `SELECT book_id, status FROM reservations WHERE id=$1 FOR UPDATE`, reservation.ID).Scan(&bookID, &mevcutStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return err
	}

	result, err := tx.ExecContext(ctx, `UPDATE reservations SET status=$1 WHERE id=$2`, reservation.Status, reservation.ID)
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

	return tx.Commit()
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

func (r *ReservationRepository) ReservationGetAll(ctx context.Context) ([]*domain.ReservationFull, error) {
	query := `SELECT r.id, r.user_id, r.book_id, r.status, r.reserved_at, r.due_date,
				u.first_name, u.last_name, b.title, b.author
			FROM reservations AS r
			INNER JOIN users AS u ON r.user_id = u.id
			INNER JOIN books AS b ON r.book_id = b.id
			ORDER BY r.reserved_at DESC`

	reservations := []*domain.ReservationFull{}
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		reserve := &domain.ReservationFull{}
		err := rows.Scan(
			&reserve.ID, &reserve.UserID, &reserve.BookID, &reserve.Status,
			&reserve.ReservedAt, &reserve.DueDate, &reserve.FirstName, &reserve.LastName,
			&reserve.BookTitle, &reserve.BookAuthor,
		)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reserve)
	}
	return reservations, nil
}

func (r *ReservationRepository) ReservationGetUserByID(ctx context.Context, userID int) ([]*domain.ReservationFull, error) {

	query := `SELECT r.id, r.user_id, r.book_id, r.status, r.reserved_at, r.due_date, 
                u.first_name, u.last_name, b.title, b.author
            FROM reservations AS r
            INNER JOIN users AS u ON r.user_id = u.id
            INNER JOIN books AS b ON r.book_id = b.id
            WHERE u.id=$1
            ORDER BY r.reserved_at DESC`

	reservations := []*domain.ReservationFull{}
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		reserve := &domain.ReservationFull{}
		err := rows.Scan(
			&reserve.ID, &reserve.UserID, &reserve.BookID, &reserve.Status,
			&reserve.ReservedAt, &reserve.DueDate, &reserve.FirstName, &reserve.LastName,
			&reserve.BookTitle, &reserve.BookAuthor,
		)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reserve)
	}
	return reservations, nil
}
