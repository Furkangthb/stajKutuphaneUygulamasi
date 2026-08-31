package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"strings"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
	"github.com/redis/go-redis/v9"
)

const booksAllIDsKey = "books:allids"

type BookRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewBookRepository(db *sql.DB, redisClient *redis.Client) *BookRepository {
	return &BookRepository{db: db, redis: redisClient}
}

func bookCacheKey(id int64) string {
	return fmt.Sprintf("book:%d", id)
}

func fetchBookFromDB(ctx context.Context, db *sql.DB, id int64) (*domain.Book, error) {
	query := `SELECT b.id, b.isbn, b.title, b.author, b.genre, b.publish_date, b.description,
			NOT EXISTS (
				SELECT 1 FROM reservations res
				WHERE res.book_id = b.id AND res.status IN ('active','pending','expired')
			) AS available
			FROM books b
			WHERE b.id=$1`

	book := &domain.Book{}
	err := db.QueryRowContext(ctx, query, id).Scan(
		&book.ID, &book.ISBN, &book.Title, &book.Author, &book.Genre, &book.PublishDate, &book.Description, &book.Available,
	)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func refreshBookCacheEntry(ctx context.Context, db *sql.DB, r *redis.Client, id int64) {
	if r == nil {
		return
	}
	book, err := fetchBookFromDB(ctx, db, id)
	if err != nil {
		return
	}
	data, err := json.Marshal(book)
	if err != nil {
		return
	}
	r.Set(ctx, bookCacheKey(id), data, 0)
	r.SAdd(ctx, booksAllIDsKey, id)
}

func removeBookCacheEntry(ctx context.Context, r *redis.Client, id int64) {
	if r == nil {
		return
	}
	r.Del(ctx, bookCacheKey(id))
	r.SRem(ctx, booksAllIDsKey, id)
}

func (r *BookRepository) WarmupCache(ctx context.Context) error {
	if r.redis == nil {
		return nil
	}
	query := `SELECT b.id, b.isbn, b.title, b.author, b.genre, b.publish_date, b.description,
			NOT EXISTS (
				SELECT 1 FROM reservations res
				WHERE res.book_id = b.id AND res.status IN ('active','pending')
			) AS available
			FROM books b
			ORDER BY b.id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	pipe := r.redis.Pipeline()
	count := 0
	for rows.Next() {
		b := &domain.Book{}
		if err := rows.Scan(&b.ID, &b.ISBN, &b.Title, &b.Author, &b.Genre, &b.PublishDate, &b.Description, &b.Available); err != nil {
			return err
		}
		data, err := json.Marshal(b)
		if err != nil {
			continue
		}
		pipe.Set(ctx, bookCacheKey(int64(b.ID)), data, 0)
		pipe.SAdd(ctx, booksAllIDsKey, b.ID)
		count++
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	slog.Info("Redis warmup kitap cache'e yuklendi")
	return nil
}

func (r *BookRepository) BookCreate(ctx context.Context, book *domain.Book) error {
	query := `INSERT INTO books (isbn,title,author,genre,publish_date,description) VALUES($1,$2,$3,$4,$5,$6) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, book.ISBN, book.Title, book.Author, book.Genre, book.PublishDate, book.Description).Scan(&book.ID)

	if err != nil {
		return err
	}
	book.Available = true
	refreshBookCacheEntry(ctx, r.db, r.redis, int64(book.ID))
	return nil
}

func (r *BookRepository) BookGetByID(ctx context.Context, id int64) (*domain.Book, error) {
	if r.redis != nil {
		if cached, err := r.redis.Get(ctx, bookCacheKey(id)).Result(); err == nil {
			var book domain.Book
			if jsonErr := json.Unmarshal([]byte(cached), &book); jsonErr == nil {
				return &book, nil
			}
		}
	}

	book, err := fetchBookFromDB(ctx, r.db, id)
	if err != nil {
		return nil, err
	}
	refreshBookCacheEntry(ctx, r.db, r.redis, id)
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
	removeBookCacheEntry(ctx, r.redis, id)
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
	refreshBookCacheEntry(ctx, r.db, r.redis, int64(book.ID))
	return nil

}

func (r *BookRepository) BookList(ctx context.Context, limit, offset int) ([]*domain.Book, error) {
	if r.redis != nil {
		if books, ok := r.bookListFromCache(ctx, limit, offset); ok {
			return books, nil
		}
	}

	query := `SELECT b.id, b.isbn, b.title, b.author, b.genre, b.publish_date, b.description,
			NOT EXISTS (
				SELECT 1 FROM reservations res
				WHERE res.book_id = b.id AND res.status IN ('active','pending')
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

	if r.redis != nil {
		for _, b := range books {
			refreshBookCacheEntry(ctx, r.db, r.redis, int64(b.ID))
		}
	}
	return books, nil
}

func (r *BookRepository) bookListFromCache(ctx context.Context, limit, offset int) ([]*domain.Book, bool) {
	idStrs, err := r.redis.SMembers(ctx, booksAllIDsKey).Result()
	if err != nil || len(idStrs) == 0 {
		return nil, false
	}

	ids := make([]int, 0, len(idStrs))
	for _, s := range idStrs {
		if n, convErr := strconv.Atoi(s); convErr == nil {
			ids = append(ids, n)
		}
	}
	sort.Ints(ids)

	start := offset
	if start > len(ids) {
		start = len(ids)
	}
	end := start + limit
	if end > len(ids) {
		end = len(ids)
	}
	pageIDs := ids[start:end]
	if len(pageIDs) == 0 {
		return []*domain.Book{}, true
	}

	keys := make([]string, len(pageIDs))
	for i, id := range pageIDs {
		keys[i] = bookCacheKey(int64(id))
	}

	vals, err := r.redis.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, false
	}

	books := make([]*domain.Book, 0, len(vals))
	for _, v := range vals {
		s, isStr := v.(string)
		if !isStr {
			return nil, false
		}
		b := &domain.Book{}
		if jsonErr := json.Unmarshal([]byte(s), b); jsonErr != nil {
			return nil, false
		}
		books = append(books, b)
	}
	return books, true
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
				WHERE res.book_id = b.id AND res.status IN ('active','pending')
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
