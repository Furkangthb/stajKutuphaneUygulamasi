package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
)

type BookServices struct {
	repo domain.BookRepository
}

func NewBookServices(repo domain.BookRepository) *BookServices {
	return &BookServices{repo: repo}
}

func (s *BookServices) BookAdd(ctx context.Context, isbn string, title string, author string, genre string, publishDate time.Time, description string) (*domain.Book, error) {
	yeniKitap := domain.Book{
		ISBN:        isbn,
		Title:       title,
		Author:      author,
		Genre:       genre,
		PublishDate: publishDate,
		Description: description,
	}
	err := s.repo.BookCreate(ctx, &yeniKitap)
	if err != nil {
		fmt.Println("KİTAP EKLENİRKEN DB HATASI:", err)
		return nil, errors.New("yeni kitap olusturalamadi")
	}
	return &yeniKitap, nil
}

func (s *BookServices) BookGet(ctx context.Context, id int64) (*domain.Book, error) {
	book, err := s.repo.BookGetByID(ctx, id)
	if err != nil {
		return nil, errors.New("kitap getirelemedi")
	}
	return book, nil
}

func (s *BookServices) BookDelete(ctx context.Context, id int64) error {
	err := s.repo.BookDelete(ctx, id)
	if err != nil {
		return errors.New("kitap silinemedi")
	}
	return nil

}

func (s *BookServices) BookUpdate(ctx context.Context, id int64, isbn string, title string, author string, genre string, publich_date time.Time, description string) (*domain.Book, error) {
	NewBook := domain.Book{
		ID:          int(id),
		ISBN:        isbn,
		Title:       title,
		Author:      author,
		Genre:       genre,
		PublishDate: publich_date,
		Description: description,
	}

	err := s.repo.BookUpdate(ctx, &NewBook)
	if err != nil {
		return nil, errors.New("kitap güncellenemedi")

	}
	return &NewBook, nil
}

func (s *BookServices) BookList(ctx context.Context, page int, page_size int) ([]*domain.Book, error) {
	if page <= 0 {
		page = 1
	}
	if page_size < 0 || page_size > 100 {
		page_size = 20
	}
	offset := (page - 1) * page_size
	book, err := s.repo.BookList(ctx, page_size, offset)
	if err != nil {
		return nil, errors.New("kitap listelemedi")
	}
	return book, nil

}

func (s *BookServices) BookSearch(ctx context.Context, keywords []string, limit int) ([]*domain.Book, error) {
	books, err := s.repo.BookSearch(ctx, limit, keywords)
	if err != nil {
		return nil, err
	}
	return books, nil

}
