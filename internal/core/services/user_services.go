package services

import (
	"context"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserServices struct {
	repo domain.UserRepository
}

func NewUserServices(repo domain.UserRepository) *UserServices {
	return &UserServices{repo: repo}
}

func (s *UserServices) UserRegister(ctx context.Context, ad string, soyad string, phone string, email string, sifre string, role string) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(sifre), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	yeniKullanici := domain.User{
		FirstName:    ad,
		LastName:     soyad,
		Phone:        phone,
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
	}
	err = s.repo.Create(ctx, &yeniKullanici)
	if err != nil {
		return nil, err
	}
	return &yeniKullanici, nil
}

func (s *UserServices) UserLogin(ctx context.Context, email string, password string) (*domain.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserServices) UserDelete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserServices) UserGet(ctx context.Context, id int64) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserServices) UserUpdate(ctx context.Context, id int, first_name string, last_name string, email string, role string) (*domain.User, error) {
	NewUser := domain.User{
		ID:        id,
		FirstName: first_name,
		LastName:  last_name,
		Email:     email,
		Role:      role,
	}
	err := s.repo.Update(ctx, &NewUser)
	if err != nil {
		return nil, err
	}
	return &NewUser, nil
}

func (s *UserServices) UserList(ctx context.Context, page int, page_size int) ([]*domain.User, error) {
	if page < 0 {
		page = 1
	}
	if page_size < 1 || page_size > 100 {
		page_size = 20
	}
	offset := (page - 1) * page_size
	users, err := s.repo.List(ctx, page_size, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}
