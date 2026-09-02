package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
)

type fakeAuthRepo struct {
	auth *domain.AuthRepository

	BlackListTokenErr   error
	IsTokenBlackListErr error

	BlackListBool bool

	BlackListDuration  time.Duration
	BlackListSignature string
}

func (f *fakeAuthRepo) BlackListToken(ctx context.Context, signature string, duration time.Duration) error {
	f.BlackListDuration = duration
	f.BlackListSignature = signature
	f.BlackListBool = true
	return f.BlackListTokenErr
}

func (f *fakeAuthRepo) IsTokenBlackList(ctx context.Context, signature string) (bool, error) {
	return f.BlackListBool, f.IsTokenBlackListErr
}

func buildToken(t *testing.T, s *AuthServices, userID int, role string, exp time.Time) string {
	t.Helper()
	payload := tokenPayload{UserID: userID, Role: role, Exp: exp.Unix()}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("payload marshal edilemedi %v", err)
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signature := s.sign(encodedPayload)
	return encodedPayload + "." + signature
}

func TestAuntSign_AyniImza(t *testing.T) {
	fakeRepo := &fakeAuthRepo{}
	s := NewAuthServices(fakeRepo, nil, "gizli")

	sign1 := s.sign("ayni")
	sign2 := s.sign("ayni")
	if sign1 != sign2 {
		t.Errorf("Ayni cikması lazimdi,farkli imzalar üretti  %s!=%s ", sign1, sign2)
	}
}

func TestAuthSign_FarkliImza(t *testing.T) {
	fakeRepo := &fakeAuthRepo{}
	s := NewAuthServices(fakeRepo, nil, "gizli")

	sign1 := s.sign("ayni")
	sign2 := s.sign("farkli")

	if sign1 == sign2 {
		t.Errorf("Farkli degerler cikmasi gerekirken,ayni degerleri dondu  %s==,%s", sign1, sign2)
	}
}

func TestAuthLogin_Basarili(t *testing.T) {
	dogruHash := mustHash(t, "123456")
	fakeRepo := &fakeUserRepo{user: &domain.User{ID: 1, Email: "furkan@gmail.com", Role: "user", PasswordHash: dogruHash}}
	s := NewAuthServices(&fakeAuthRepo{}, fakeRepo, "gizli")
	token, err := s.Login(context.Background(), "furkan@gmail.com", "123456")
	if err != nil {
		t.Fatalf("beklenmeyen hata %v", err)
	}
	parcalar := strings.Split(token, ".")
	if len(parcalar) != 2 {
		t.Errorf("Token 2 parcali olmali farkli geldi %d", len(parcalar))
	}

}

func TestAuthLogin_BosDondu_Hata(t *testing.T) {
	fakeRepo := &fakeUserRepo{GetByEmailErr: errTest("hata")}
	s := NewAuthServices(&fakeAuthRepo{}, fakeRepo, "gizli")
	_, err := s.Login(context.Background(), "furkan@gmail.com", "123456")
	if err == nil {
		t.Fatalf("Hata bekleniyordu,nil dondu")
	}
}

func TestAuthLogin_YanlisSifre_Hata(t *testing.T) {
	dogruHash := mustHash(t, "123546")
	fakeRepo := &fakeUserRepo{user: &domain.User{ID: 1, Email: "furkan@gmail.com", Role: "user", PasswordHash: dogruHash}}
	s := NewAuthServices(&fakeAuthRepo{}, fakeRepo, "gizli")
	_, err := s.Login(context.Background(), "furkan@gmail.com", "147258")
	if err == nil {
		t.Fatalf("Hata bekleniyordu ,nil dondu")
	}
}

func TestAuthVerifyToken_Basarali(t *testing.T) {
	s := NewAuthServices(&fakeAuthRepo{}, nil, "gizli")
	token := buildToken(t, s, 1, "user", time.Now().Add(24*time.Hour))
	userID, role, err := s.VerifyToken(context.Background(), token)
	if err != nil {
		t.Fatalf("Beklenmeyen hata %v", err)
	}
	if role != "user" {
		t.Errorf("kullanici role yanlis %v", role)
	}
	if userID != 1 {
		t.Errorf("ID ler eslesmiyor ,beklenen=1  gelen=%d", userID)
	}

}

func TestAuthVerifyToken_BozukFormat_Hata(t *testing.T) {
	s := NewAuthServices(&fakeAuthRepo{}, nil, "gizli")
	_, _, err := s.VerifyToken(context.Background(), "yanlis.formattaki.bir.tane.veri")

	if err == nil {
		t.Fatalf("Hata bekleniyordu,nil dondu")
	}
}

func TestAuthVerifyToken_KurcalanmisToken_Hata(t *testing.T) {

	s := NewAuthServices(&fakeAuthRepo{}, nil, "gizli")
	token := buildToken(t, s, 5, "user", time.Now().Add(24*time.Hour))

	parcalar := strings.Split(token, ".")
	kurcalanmisToken := parcalar[0] + "." + parcalar[1] + "2"

	_, _, err := s.VerifyToken(context.Background(), kurcalanmisToken)
	if err == nil {
		t.Fatalf("Hata bekleniyordu ,nil dondu")
	}

}

func TestAuthVerifyToken_TokenSureDolmus_Hata(t *testing.T) {
	s := NewAuthServices(&fakeAuthRepo{}, nil, "gizli")
	suresiDolmusToken := buildToken(t, s, 7, "admin", time.Now().Add(-1*time.Hour))
	_, _, err := s.VerifyToken(context.Background(), suresiDolmusToken)
	if err == nil {
		t.Fatalf("Hata bekleniyordu,nil dondu")
	}
}

func TestAuthVerifyToken_Blacklistte_Hata(t *testing.T) {
	fakeRepo := &fakeAuthRepo{BlackListBool: true}
	s := NewAuthServices(fakeRepo, nil, "gizli")

	token := buildToken(t, s, 5, "admin", time.Now().Add(24*time.Hour))
	_, _, err := s.VerifyToken(context.Background(), token)

	if err == nil {
		t.Fatalf("Hata bekleniyordu,nil dondu")
	}
}

func TestAuthLogout_DogruImzaylaBlacklisteEkler_Basarali(t *testing.T) {
	fakeRepo := &fakeAuthRepo{}
	s := NewAuthServices(fakeRepo, nil, "gizli")
	token := buildToken(t, s, 5, "user", time.Now().Add(24*time.Hour))
	beklenenImza := strings.Split(token, ".")[1]
	err := s.Logout(context.Background(), token)
	if err != nil {
		t.Fatalf("Beklenmeyen hata %v", err)
	}
	if !fakeRepo.BlackListBool {
		t.Fatalf("Blacklist cagrilamadi")
	}
	if fakeRepo.BlackListSignature != beklenenImza {
		t.Errorf("imzalar uyusmadi , beklenen=%q  gelen=%q", beklenenImza, fakeRepo.BlackListSignature)
	}
	if fakeRepo.BlackListDuration <= 0 {
		t.Errorf("Blacklistteki zaman eksi deger almamali %v geldi", fakeRepo.BlackListDuration)
	}
}

func TestAuthLogout_KurcalanmisToken_Hata(t *testing.T) {

	s := NewAuthServices(&fakeAuthRepo{}, nil, "gizli")
	token := buildToken(t, s, 5, "user", time.Now().Add(24*time.Hour))

	parcalar := strings.Split(token, ".")
	kurcalanmisToken := parcalar[0] + "." + parcalar[1] + "2"

	_, _, err := s.VerifyToken(context.Background(), kurcalanmisToken)
	if err == nil {
		t.Fatalf("Hata bekleniyordu ,nil dondu")
	}

}

func TestAuthLogout_TokenSureDolmus_Hata(t *testing.T) {
	s := NewAuthServices(&fakeAuthRepo{}, nil, "gizli")
	suresiDolmusToken := buildToken(t, s, 7, "admin", time.Now().Add(-1*time.Hour))
	_, _, err := s.VerifyToken(context.Background(), suresiDolmusToken)
	if err == nil {
		t.Fatalf("Hata bekleniyordu,nil dondu")
	}
}
