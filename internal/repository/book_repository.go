package repository

import (
	"context"
	"database/sql"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) BookCreate(ctx context.Context, book *domain.Book) error {
	query := `INSERT INTO books (title,author,genre,publish_date,description,stock_count) VALUES($1,$2,$3,$4,$5,$6) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, book.Title, book.Author, book.Genre, book.PublishDate, book.Description, book.StockCount).Scan(&book.ID)

	if err != nil {
		return nil
	}
	return nil
}

func (r *BookRepository) BookGetByID(ctx context.Context, id int64) (*domain.Book, error) {
	query := `SELECT id,title,author,genre,publish_date,description,stock_count
			FROM book
			WHERE id=$1`

	book := &domain.Book{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&book.ID, &book.Title, &book.Author, &book.Genre, &book.PublishDate, &book.Description, &book.StockCount)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (r *BookRepository) BookDelete(ctx context.Context, id int64) error {
	query := `DELETE 
			FROM book
			WHERE id=&1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return err
	}
	return nil
}

func (r *BookRepository) BookUpdate(ctx context.Context, book *domain.Book) error {
	query := `UPDATE books
			SET title=$1,author=$2,genre=$3,publish_date=$4,description=$5,stock_count=$6
			WHERE id=$7`

	result, err := r.db.ExecContext(ctx, query, book.Title, book.Author, book.Genre, book.PublishDate, book.Description, book.StockCount)
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
	query := `SELECT title,author,genre,publish_date,description,stock_count
			FROM books
			ORDER BY id
			LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*domain.Book
	for rows.Next() {
		b := &domain.Book{}
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Genre, &b.PublishDate, &b.Description, &b.StockCount); err != nil {
			return nil, err
		}
		books = append(books, b)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}
