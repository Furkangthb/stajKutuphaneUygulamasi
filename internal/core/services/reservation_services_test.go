package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/repository"
)

type fakeReservationRepo struct { // İsmindeki 'a' harfini de 'e' olarak düzelttim :)
	createErr    error
	updateErr    error
	getIdErr     error
	getAllErr    error
	getUserIdErr error
	expireErr    error

	reservation     *domain.Reservation
	reservationList []*domain.ReservationFull
	expiredIDs      []int

	createdReservation *domain.Reservation
	updatedReservation *domain.Reservation
	getUserCalledWith  int
}

func (f *fakeReservationRepo) ReservationCreate(ctx context.Context, r *domain.Reservation) error {
	f.createdReservation = r
	if f.createErr != nil {
		return f.createErr
	}
	r.ID = 1
	r.ReservedAt = time.Now()
	return nil
}

func (f *fakeReservationRepo) ReservationUpdate(ctx context.Context, r *domain.Reservation) error {
	f.updatedReservation = r
	return f.updateErr
}

// DİKKAT: Id yerine ID yazıldı!
func (f *fakeReservationRepo) ReservationGetByID(ctx context.Context, id int) (*domain.Reservation, error) {
	if f.getIdErr != nil {
		return nil, f.getIdErr
	}
	return f.reservation, nil
}

func (f *fakeReservationRepo) ReservationGetAll(ctx context.Context) ([]*domain.ReservationFull, error) {
	if f.getAllErr != nil {
		return nil, f.getAllErr
	}
	return f.reservationList, nil
}

func (f *fakeReservationRepo) ReservationGetUserByID(ctx context.Context, userID int) ([]*domain.ReservationFull, error) {
	f.getUserCalledWith = userID
	if f.getUserIdErr != nil {
		return nil, f.getUserIdErr
	}
	return f.reservationList, nil
}

func (f *fakeReservationRepo) ExpireOverdue(ctx context.Context) ([]int, error) {
	return f.expiredIDs, f.expireErr
}

func TestReservationCreate_Basarali(t *testing.T) {
	fakeRepo := &fakeReservationRepo{}
	s := NewReservationServices(fakeRepo)

	reservation, err := s.ReservationCreate(context.Background(), 2, "user", 3, 4)
	if err != nil {
		t.Fatalf("beklenmeyen bir hata olustu  %v", err)
	}

	if reservation.BookID != 4 {
		t.Errorf("id ler eslesmiyor , beklenen=4 gelen=%d", reservation.BookID)
	}
	if fakeRepo.createdReservation == nil {
		t.Fatalf("repo hic cagırılmamıs")
	}
}

func TestReservationCreate_KitapBulunamadi_Hata(t *testing.T) {
	fakeRepo := &fakeReservationRepo{createErr: repository.ErrBookNotFound}
	s := NewReservationServices(fakeRepo)

	_, err := s.ReservationCreate(context.Background(), 1, "user", 5, 4)
	if err == nil {
		t.Fatal("Hata bekleniyordu,nil geldi")
	}
}

func TestReservaitonCreate_KitapMusaitDegil_Hata(t *testing.T) {
	fakeRepo := &fakeReservationRepo{createErr: repository.ErrBookNotAvailable}
	s := NewReservationServices(fakeRepo)

	_, err := s.ReservationCreate(context.Background(), 1, "user", 5, 4)
	if err == nil {
		t.Fatalf("Hata bekleniyordu,nil dondu")
	}

}

func TestReservationCreate_LimitAsildi_Hata(t *testing.T) {
	fakeRepo := &fakeReservationRepo{createErr: repository.ErrReservationLimitExceeded}
	s := NewReservationServices(fakeRepo)

	_, err := s.ReservationCreate(context.Background(), 1, "user", 5, 4)
	if err == nil {
		t.Fatalf("Hata bekleniyordu,nil dondu")
	}
}

func TestReservationUpdate_Basarali(t *testing.T) {
	fakeRepo := &fakeReservationRepo{reservation: &domain.Reservation{ID: 7, UserID: 3, Status: domain.StatusPending}}
	s := NewReservationServices(fakeRepo)
	reservation, err := s.ReservationUpdate(context.Background(), 3, "admin", 7, domain.StatusActive)
	if err != nil {
		t.Fatalf("Beklenmeyen bir hata olustu %v", err)
	}
	if reservation.Status != domain.StatusActive {
		t.Errorf("durumlar eslesmiyor")
	}

}

func TestReservationUpdate_OlmayanDurum_Hata(t *testing.T) {
	fakeRepo := &fakeReservationRepo{}
	s := NewReservationServices(fakeRepo)

	_, err := s.ReservationUpdate(context.Background(), 3, "admin", 7, "olmayanDurum")
	if err == nil {
		t.Fatalf("Hata bekleniyordu,nil dondu %v", err)
	}

}

func TestReservationUpdate_ReservationBulunamadi_Hata(t *testing.T) {
	fakeRepo := &fakeReservationRepo{}
	s := NewReservationServices(fakeRepo)

	_, err := s.ReservationUpdate(context.Background(), 3, "admin", 7, domain.StatusActive)
	if !errors.Is(err, ErrReservationNotFound) {
		t.Errorf("Hata bekleniyordu,nil dondu %v", err)
	}
}

func TestReservationUpdate_YetkisizReservation_Hata(t *testing.T) {
	fakeRepo := &fakeReservationRepo{reservation: &domain.Reservation{ID: 1, UserID: 5, Status: domain.StatusPending}}
	s := NewReservationServices(fakeRepo)

	_, err := s.ReservationUpdate(context.Background(), 1, "user", 10, domain.StatusCancelled)
	if !errors.Is(err, ErrForbidden) {
		t.Errorf("Hata bekleniyordu,nil dondu  %v", err)
	}
}

func TestReservationUpdate_DurumGecis_Hata(t *testing.T) {
	fakeRepo := &fakeReservationRepo{reservation: &domain.Reservation{ID: 1, UserID: 5, Status: domain.StatusCompleted}}
	s := NewReservationServices(fakeRepo)

	_, err := s.ReservationUpdate(context.Background(), 5, "user", 10, domain.StatusActive)

	if !errors.Is(err, ErrInvalidTransition) {
		t.Errorf("Yanlis hata dondu. Beklenen: %v, Gelen: %v", ErrInvalidTransition, err)
	}
}

func TestReservationGetById_Basarali(t *testing.T) {
	fakeRepo := &fakeReservationRepo{reservation: &domain.Reservation{ID: 10}}
	s := NewReservationServices(fakeRepo)
	_, err := s.ReservationGetByID(context.Background(), 10)
	if err != nil {
		t.Fatalf("Beklenmeyen bir hata olustu %v", err)
	}
}

func TestReservationGetById_Hata(t *testing.T) {
	fakeRepo := &fakeReservationRepo{getIdErr: errTest("hata")}
	s := NewReservationServices(fakeRepo)
	_, err := s.ReservationGetByID(context.Background(), 10)
	if err == nil {
		t.Fatalf("Hata bekleniyordu,nil dondu")
	}
}

func TestReservationListAll_Basarali(t *testing.T) {
	fakeRepo := &fakeReservationRepo{reservationList: []*domain.ReservationFull{{ID: 1}, {ID: 2}}}
	s := NewReservationServices(fakeRepo)
	list, err := s.ReservationListAll(context.Background())
	if err != nil {
		t.Fatalf("Beklenmeyen bir hata olustu %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("beklenen degerden farklı bir sonuc geldi   %d", len(list))
	}
}

func TestReservationListAll_Hata(t *testing.T) {
	fakeRepo := &fakeReservationRepo{getAllErr: errTest("hata")}
	s := NewReservationServices(fakeRepo)
	_, err := s.ReservationListAll(context.Background())
	if err == nil {
		t.Fatalf("Hata bekleniyordu,nil dondu")

	}
}

func TestReservationGetByUserID_Basarali(t *testing.T) {
	fakeRepo := &fakeReservationRepo{reservationList: []*domain.ReservationFull{{ID: 1}, {ID: 2}}}
	s := NewReservationServices(fakeRepo)

	list, err := s.ReservationGetByUserID(context.Background(), 3)
	if err != nil {
		t.Fatalf("Beklenmeyen bir hata olustu  %v", err)
	}

	if len(list) != 2 {
		t.Fatalf("2 deger donmesi bekleniyordu ama %d dondu", len(list))
	}

	if fakeRepo.getUserCalledWith != 3 {
		t.Errorf("id ler eslesmiyor , beklenen=3 gelen=%v", fakeRepo.getUserCalledWith)
	}
}

func TestReservationGetByUserID_Hata(t *testing.T) {
	fakeRepo := &fakeReservationRepo{getUserIdErr: errTest("hata")}
	s := NewReservationServices(fakeRepo)
	_, err := s.ReservationGetByUserID(context.Background(), 3)
	if err == nil {
		t.Fatalf("Hata bekleniyordu,nil")
	}
}
