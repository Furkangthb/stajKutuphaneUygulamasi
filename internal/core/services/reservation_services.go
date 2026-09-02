package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/repository"
)

var (
	ErrReservationNotFound = errors.New("rezervasyon bulunamadi")
	ErrForbidden           = errors.New("bu islem icin yetkiniz yok")
	ErrInvalidStatus       = errors.New("gecersiz durum degeri")
	ErrInvalidTransition   = errors.New("bu durum gecisine izin verilmiyor")
)

type ReservationServices struct {
	repo domain.ReservationRepository
}

func NewReservationServices(repo domain.ReservationRepository) *ReservationServices {
	return &ReservationServices{repo: repo}
}

func (s *ReservationServices) ReservationCreate(ctx context.Context, requesterID int, requesterRole string, bodyUserID int, bookId int) (*domain.Reservation, error) {
	targetUserID := requesterID
	if requesterRole == "admin" && bodyUserID != 0 {
		targetUserID = bodyUserID
	}

	dueDate := time.Now().AddDate(0, 0, 14)
	newReservation := domain.Reservation{
		UserID:  targetUserID,
		BookID:  bookId,
		Status:  domain.StatusPending,
		DueDate: dueDate,
	}
	err := s.repo.ReservationCreate(ctx, &newReservation)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrBookNotFound):
			return nil, errors.New("kitap bulunamadi")
		case errors.Is(err, repository.ErrBookNotAvailable):
			return nil, errors.New("bu kitap su anda musait degil (zaten rezerve/odunc)")
		case errors.Is(err, repository.ErrReservationLimitExceeded):
			return nil, fmt.Errorf("en fazla %d acik rezervasyonunuz olabilir, once mevcutlardan birini iade/iptal edin", domain.MaxActiveReservationsPerUser)
		default:
			slog.Error("Reservasyon db hatasi", slog.Any("error", err))

			return nil, errors.New("rezervasyon yapilamadi")
		}
	}
	return &newReservation, nil
}

func (s *ReservationServices) ReservationUpdate(ctx context.Context, requesterID int, requesterRole string, id int, newStatus string) (*domain.Reservation, error) {
	if !domain.IsValidStatus(newStatus) {
		return nil, ErrInvalidStatus
	}

	existing, err := s.repo.ReservationGetByID(ctx, id)
	if err != nil || existing == nil {
		return nil, ErrReservationNotFound
	}

	if requesterRole != "admin" && existing.UserID != requesterID {
		return nil, ErrForbidden
	}

	if !domain.CanUserTransition(requesterRole, existing.Status, newStatus) {
		return nil, ErrInvalidTransition
	}

	reworkReservation := domain.Reservation{
		ID:     id,
		Status: newStatus,
	}
	err = s.repo.ReservationUpdate(ctx, &reworkReservation)
	if err != nil {
		return nil, errors.New("rezervasyon guncellenmedi")
	}
	return &reworkReservation, nil
}

func (s *ReservationServices) ReservationGetByID(ctx context.Context, id int) (*domain.Reservation, error) {
	reservation, err := s.repo.ReservationGetByID(ctx, id)
	if err != nil || reservation == nil {
		return nil, errors.New("rezervasyon getirilemedi veya bulunamadi")
	}
	return reservation, nil
}

func (s *ReservationServices) ReservationListAll(ctx context.Context) ([]*domain.ReservationFull, error) {
	reservations, err := s.repo.ReservationGetAll(ctx)
	if err != nil {
		return nil, errors.New("rezervasyonlar getirilemedi")
	}
	return reservations, nil
}

func (s *ReservationServices) ReservationGetByUserID(ctx context.Context, id int) ([]*domain.ReservationFull, error) {
	reservations, err := s.repo.ReservationGetUserByID(ctx, id)
	if err != nil {
		return nil, errors.New("reservasyonlar getiremedi")
	}
	return reservations, nil
}

func (s *ReservationServices) ExpireOverdueReservations(ctx context.Context) (int, error) {
	bookIDs, err := s.repo.ExpireOverdue(ctx)
	if err != nil {
		return 0, err
	}
	return len(bookIDs), nil
}
