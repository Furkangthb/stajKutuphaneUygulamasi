package services

import (
	"context"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserServices struct {
	repo *repository.UserRepository
}

func NewUserServices(repo *repository.UserRepository) *UserServices {
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
	err = s.repo.UserRegister(ctx, &yeniKullanici)
	if err != nil {
		return nil, err
	}
	return &yeniKullanici, nil
}
