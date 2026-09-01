package services

import (
	"context"
	"testing"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepo struct {
	user *domain.User

	CreateErr         error
	DeleteErr         error
	UpdateErr         error
	ListErr           error
	GetByIDErr        error
	GetByEmailErr     error
	UpdatePasswordErr error

	CreateUser        *domain.User
	UpdatedUser       *domain.User
	listLimit         int
	listOffset        int
	updatedPasswordID int64
	updatedHash       string
}

func (f *fakeUserRepo) Create(ctx context.Context, u *domain.User) error {
	f.CreateUser = u
	return f.CreateErr
}

func (f *fakeUserRepo) Delete(ctx context.Context, id int64) error {
	return f.DeleteErr
}

func (f *fakeUserRepo) Update(ctx context.Context, u *domain.User) error {
	f.UpdatedUser = u
	return f.UpdateErr
}

func (f *fakeUserRepo) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	f.listLimit = limit
	f.listOffset = offset
	return nil, f.ListErr
}

func (f *fakeUserRepo) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	if f.GetByIDErr != nil {
		return nil, f.GetByIDErr
	}
	return f.user, nil
}

func (f *fakeUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if f.GetByEmailErr != nil {
		return nil, f.GetByEmailErr
	}
	return f.user, nil
}

func (f *fakeUserRepo) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	f.updatedPasswordID = id
	f.updatedHash = passwordHash
	return f.UpdatePasswordErr
}

func mustHash(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash uretilemedi: %v", err)
	}
	return string(hash)
}

func TestUserRegister_Basarili(t *testing.T) {
	fakeRepo := &fakeUserRepo{}
	s := NewUserServices(fakeRepo)
	user, err := s.UserRegister(context.Background(), "Furkan", "Yüksel", "12345678911", "furkan@gmail.com", "123456")
	if err != nil {
		t.Fatalf("UserRegister basarisiz oldu: %v", err)
	}
	if user.FirstName != "Furkan" || user.LastName != "Yüksel" {
		t.Errorf("kullanici adi veya soyismi yanlış =%v", user)
	}
	if user.Email != "furkan@gmail.com" {
		t.Errorf("kullanici emaili yanlış =%v", user)
	}
	if user.Phone != "12345678911" {
		t.Errorf("kullanici telefonu yanlış =%v", user)
	}
	if user.Role != "user" {
		t.Errorf("yanlış role donusu =%v", user.Role)
	}
	if user.PasswordHash == "123456" {
		t.Errorf("sifre hashlenemden saklanıyor =%v", user.PasswordHash)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("123456")); err != nil {
		t.Errorf("hash ile hashlenmis sifre eslenmiyor")
	}
	if fakeRepo.CreateUser == nil {
		t.Fatal("kullanıcı olusturulamiyor")
	}
}

func TestUserRegister_Hata(t *testing.T) {
	fakeRepo := &fakeUserRepo{CreateErr: errTest("hata")}
	s := NewUserServices(fakeRepo)
	_, err := s.UserRegister(context.Background(), "Furkan", "Yüksel", "12345678911", "furkan@gmail.com", "123456")
	if err == nil {
		t.Fatal("hata bekliyorduk ama olmadı")
	}
}

func TestUserLogin_Basarili(t *testing.T) {
	dogruHash := mustHash(t, "123456")
	user := &domain.User{ID: 5, Email: "furkan@gmail.com", PasswordHash: dogruHash}
	fakeRepo := &fakeUserRepo{user: user}
	s := NewUserServices(fakeRepo)

	user, err := s.UserLogin(context.Background(), "furkan@gmail.com", "123456")
	if err != nil {
		t.Fatalf(" UserLogin basarisiz oldu: %v", err)
	}
	if user.ID != 5 {
		t.Errorf("yanlış kullanıcı ID'si =%v", user.ID)
	}

}

func TestUserLogin_Hata(t *testing.T) {
	fakeRepo := &fakeUserRepo{GetByEmailErr: errTest("hata")}
	s := NewUserServices(fakeRepo)
	_, err := s.UserLogin(context.Background(), "furkan@gmail.com", "123456")
	if err == nil {
		t.Fatal("hata bekliyorduk ama olmadı")
	}
}

func TestUserDelete_Basarili(t *testing.T) {
	fakeRepo := &fakeUserRepo{}
	s := NewUserServices(fakeRepo)
	err := s.UserDelete(context.Background(), 5)
	if err != nil {
		t.Fatalf("beklenmeyen hata oluştu: %v", err)
	}
}

func TestUserDelete_Hata(t *testing.T) {
	fakeRepo := &fakeUserRepo{DeleteErr: errTest("hata")}
	s := NewUserServices(fakeRepo)
	err := s.UserDelete(context.Background(), 5)
	if err == nil {
		t.Fatalf("hata bekleniyordu,nul geldi")
	}
}

func TestUserLogin_SifreHatasi(t *testing.T) {
	dogruHash := mustHash(t, "123456")
	fakeRepo := &fakeUserRepo{user: &domain.User{ID: 5, Email: "furkan@gmail.com", PasswordHash: dogruHash}}

	s := NewUserServices(fakeRepo)
	_, err := s.UserLogin(context.Background(), "furkan@gmail.com", "wrongpassword")
	if err == nil {
		t.Fatal("hata bekliyorduk ama olmadı")
	}
}

func TestUserGet_Basarili(t *testing.T) {
	user := &domain.User{ID: 5, Email: "furkan@gmail.com"}
	fakeRepo := &fakeUserRepo{user: user}
	s := NewUserServices(fakeRepo)
	user, err := s.UserGet(context.Background(), 5)
	if err != nil {
		t.Fatalf("beklenmeyen hata olustu %v ", err)
	}

	if user.ID != 5 {
		t.Errorf("id ler eslesmiyor  %v", user.ID)
	}
}

func TestUserGet_Hata(t *testing.T) {
	fakeRepo := &fakeUserRepo{GetByIDErr: errTest("hata")}
	s := NewUserServices(fakeRepo)
	_, err := s.UserGet(context.Background(), 3)
	if err == nil {
		t.Fatalf("hata bekleniyordu,nil dondu")
	}
}

func TestUserUpdate_Basarili(t *testing.T) {
	fakeRepo := &fakeUserRepo{}
	s := NewUserServices(fakeRepo)
	user, err := s.UserUpdate(context.Background(), 5, "Furkan", "Yüksel", "furkan123@gmail.com", "admin")
	if err != nil {
		t.Fatalf("beklenmeyen hata olustu %v  ", err)
	}
	if user.FirstName != "Furkan" || user.LastName != "Yüksel" || user.Email != "furkan123@gmail.com" || user.Role != "admin" {
		t.Errorf("guncellenen alanlar yanlis dondu =%v", user)
	}
	if fakeRepo.UpdatedUser == nil {
		t.Fatalf("guncellenen kullanıcı repo ya gitmedi")
	}
	if fakeRepo.UpdatedUser.PasswordHash != "" {
		t.Errorf("sifre yanlısıkla sıfırlanabilir")
	}
}

func TestUserUpdate_Hata(t *testing.T) {
	fakeRepo := &fakeUserRepo{UpdateErr: errTest("hata")}
	s := NewUserServices(fakeRepo)
	_, err := s.UserUpdate(context.Background(), 5, "Furkan", "Yüksel", "furkan123@gmail.com", "admin")
	if err == nil {
		t.Fatalf("hata bekleniyordu ,nil dondu")
	}
}

func TestUserList_Basarili(t *testing.T) {
	fakeRepo := &fakeUserRepo{}
	s := NewUserServices(fakeRepo)
	_, err := s.UserList(context.Background(), 0, 5)
	if err != nil {
		t.Fatalf("beklenmeyen hata olustu %v ", err)
	}
	if fakeRepo.listLimit != 5 || fakeRepo.listOffset != 0 {
		t.Errorf("beklenmeyen limit veya offset! Beklenen: L=5, O=0 | Gelen: Limit=%v, Offset=%v", fakeRepo.listLimit, fakeRepo.listOffset)
	}

}

func TestUserList_Hata(t *testing.T) {
	fakeRepo := &fakeUserRepo{ListErr: errTest("hata")}
	s := NewUserServices(fakeRepo)
	_, err := s.UserList(context.Background(), 0, 5)
	if err == nil {
		t.Fatalf("hata bekleniyordu,nil dondu")
	}
}

func FuzzUserList_Sayfalama(f *testing.F) {
	f.Add(1, 10)
	f.Add(0, 5)

	f.Fuzz(func(t *testing.T, page int, pageSize int) {

		fakeRepo := &fakeUserRepo{}
		s := NewUserServices(fakeRepo)

		_, _ = s.UserList(context.Background(), page, pageSize)

		if fakeRepo.listOffset < 0 {
			t.Errorf("Fuzzer bir açık buldu! Negatif offset oluştu. Page: %v, Size: %v, Olusan Offset: %v",
				page, pageSize, fakeRepo.listOffset)
		}
	})
}

func TestUserChangePassword_Basarili(t *testing.T) {
	dogruHash := mustHash(t, "123456")
	fakeRepo := &fakeUserRepo{user: &domain.User{ID: 1, PasswordHash: dogruHash}}
	s := NewUserServices(fakeRepo)
	err := s.UserChangePassword(context.Background(), 1, "123456", "123456789")
	if err != nil {
		t.Fatalf("beklenmeyen hata olustu %v", err)
	}

	if fakeRepo.updatedPasswordID != 1 {
		t.Errorf("idler eslesmiyor")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(fakeRepo.updatedHash), []byte("123456789")); err != nil {
		t.Errorf("hashlenmis yeni sifre ile hashli sifre eslesmiyor")
	}
}

func TestUserChangePassword_Hata(t *testing.T) {
	dogruHash := mustHash(t, "123456")
	fakeRepo := &fakeUserRepo{user: &domain.User{ID: 1, PasswordHash: dogruHash}, UpdatePasswordErr: errTest("hata")}
	s := NewUserServices(fakeRepo)
	err := s.UserChangePassword(context.Background(), 1, "123456", "123456789")
	if err == nil {
		t.Fatalf("hata bekleniyordu ,nil dondu")
	}
}

type simpleError string

func (e simpleError) Error() string {
	return string(e)
}
func errTest(msg string) error {
	return simpleError(msg)
}
