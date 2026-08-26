package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) BookCreate(ctx context.Context, book *domain.Book) error {
	query := `INSERT INTO books (isbn,title,author,genre,publish_date,description) VALUES($1,$2,$3,$4,$5,$6) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, book.ISBN, book.Title, book.Author, book.Genre, book.PublishDate, book.Description).Scan(&book.ID)

	if err != nil {
		return err
	}
	book.Available = true
	return nil
}

func (r *BookRepository) BookGetByID(ctx context.Context, id int64) (*domain.Book, error) {
	query := `SELECT b.id, b.isbn, b.title, b.author, b.genre, b.publish_date, b.description,
			NOT EXISTS (
				SELECT 1 FROM reservations res
				WHERE res.book_id = b.id AND res.status='active'
			) AS available
			FROM books b
			WHERE b.id=$1`

	book := &domain.Book{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&book.ID, &book.ISBN, &book.Title, &book.Author, &book.Genre, &book.PublishDate, &book.Description, &book.Available,
	)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (r *BookRepository) BookDelete(ctx context.Context, id int64) error {
	query := `DELETE 
			FROM books
			WHERE id=$1`
	result, err := r.db.ExecContext(ctx, query, id)
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

func (r *BookRepository) BookUpdate(ctx context.Context, book *domain.Book) error {
	query := `UPDATE books
			SET isbn=$1,title=$2,author=$3,genre=$4,publish_date=$5,description=$6
			WHERE id=$7`

	result, err := r.db.ExecContext(ctx, query, book.ISBN, book.Title, book.Author, book.Genre, book.PublishDate, book.Description, book.ID)
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

func (r *BookRepository) BookList(ctx context.Context, limit, offset int) ([]*domain.Book, error) {
	query := `SELECT b.id, b.isbn, b.title, b.author, b.genre, b.publish_date, b.description,
			NOT EXISTS (
				SELECT 1 FROM reservations res
				WHERE res.book_id = b.id AND res.status='active'
			) AS available
			FROM books b
			ORDER BY b.id
			LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*domain.Book
	for rows.Next() {
		b := &domain.Book{}
		if err := rows.Scan(&b.ID, &b.ISBN, &b.Title, &b.Author, &b.Genre, &b.PublishDate, &b.Description, &b.Available); err != nil {
			return nil, err
		}
		books = append(books, b)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

func (r *BookRepository) BookSearch(ctx context.Context, limit int, keywords []string) ([]*domain.Book, error) {
	if len(keywords) == 0 {
		return []*domain.Book{}, nil
	}
	var conditions []string
	var args []interface{}
	var argsIndex = 1

	for _, kw := range keywords {
		conditions = append(conditions, fmt.Sprintf("(b.title ILIKE $%d OR b.author ILIKE $%d OR b.genre ILIKE $%d)", argsIndex, argsIndex, argsIndex))
		args = append(args, "%"+kw+"%")
		argsIndex++
	}

	query := fmt.Sprintf(`SELECT b.id, b.isbn, b.title, b.author, b.genre, b.publish_date, b.description,
			NOT EXISTS (
				SELECT 1 FROM reservations res
				WHERE res.book_id = b.id AND res.status='active'
			) AS available
			FROM books b
			WHERE %s
			LIMIT $%d`, strings.Join(conditions, " OR "), argsIndex)
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*domain.Book
	for rows.Next() {
		b := &domain.Book{}
		if err := rows.Scan(&b.ID, &b.ISBN, &b.Title, &b.Author, &b.Genre, &b.PublishDate, &b.Description, &b.Available); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}
