package services

import (
	"context"
	"time"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/repository"
)

type BookServices struct {
	repo *repository.BookRepository
}

func NewBookServices(repo *repository.BookRepository) *BookServices {
	return &BookServices{repo: repo}
}

func (s *BookServices) BookAdd(ctx context.Context, title string, author string, genre string, publishDate time.Time, description string, stockCount int) (*domain.Book, error) {
	yeniKitap := domain.Book{
		Title:       title,
		Author:      author,
		Genre:       genre,
		PublishDate: publishDate,
		Description: description,
		StockCount:  stockCount,
	}
	err := s.repo.BookCreate(ctx, &yeniKitap)
	if err != nil {
		return nil, err
	}
	return &yeniKitap, nil
}

func (s *BookServices) BookGet(ctx context.Context, id int64) (*domain.Book, error) {
	book, err := s.repo.BookGetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (s *BookServices) BookDelete(ctx context.Context, id int64) error {
	err := s.repo.BookDelete(ctx, id)
	if err != nil {
		return err
	}
	return err

}

func (s *BookServices) BookUpdate(ctx context.Context, id int, title string, author string, genre string, publich_date time.Time, description string, stock_count int) (*domain.Book, error) {
	NewBook := domain.Book{
		ID:          id,
		Title:       title,
		Author:      author,
		Genre:       genre,
		PublishDate: publich_date,
		Description: description,
		StockCount:  stock_count,
	}

	err := s.repo.BookUpdate(ctx, &NewBook)
	if err != nil {
		return nil, err

	}
	return &NewBook, nil
}
