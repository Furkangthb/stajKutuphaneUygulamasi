package services

import (
	"context"
	"errors"
	"time"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/repository"
)

type ReservationServices struct {
	repo *repository.ReservationRepository
}

func NewReservationServices(repo *repository.ReservationRepository) *ReservationServices {
	return &ReservationServices{repo: repo}
}
func (s *ReservationServices) ReservationCreate(ctx context.Context, userId int, bookId int, status string) (*domain.Reservation, error) {

	dueDate := time.Now().AddDate(0, 0, 14)
	newReservation := domain.Reservation{
		UserID:  userId,
		BookID:  bookId,
		Status:  status,
		DueDate: dueDate,
	}
	err := s.repo.ReservationCreate(ctx, &newReservation)
	if err != nil {
		return nil, errors.New("rezervasyon yapilamadi")
	}
	return &newReservation, nil
}

func (s *ReservationServices) ReservationUpdate(ctx context.Context, id int, status string) (*domain.Reservation, error) {
	ReworkReservation := domain.Reservation{

		ID:     id,
		Status: status,
	}
	err := s.repo.ReservationUpdate(ctx, &ReworkReservation)
	if err != nil {
		return nil, errors.New("rezervasyom guncellenmedi")
	}
	return &ReworkReservation, nil
}

func (s *ReservationServices) ReservationGetByID(ctx context.Context, id int) (*domain.Reservation, error) {
	reservation, err := s.repo.ReservationGetByID(ctx, id)
	if err != nil {
		return nil, errors.New("rezervasyon getirilemedi")
	}
	return reservation, nil
}

func (s *ReservationServices) ReservationGetByUserID(ctx context.Context, id int) ([]*domain.ReservationWithUser, error) {
	reservations, err := s.repo.ReservationGetUserByID(ctx, id)
	if err != nil {
		return nil, errors.New("reservasyonlar getiremedi")
	}
	return reservations, nil
}
