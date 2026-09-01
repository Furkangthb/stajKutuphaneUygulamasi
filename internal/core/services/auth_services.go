package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthServices struct {
	authRepo  domain.AuthRepository
	userRepo  domain.UserRepository
	secretKey []byte
}

func NewAuthServices(authRepo domain.AuthRepository, userRepo domain.UserRepository, secret string) *AuthServices {
	return &AuthServices{authRepo: authRepo, userRepo: userRepo, secretKey: []byte(secret)}
}

type tokenPayload struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
}

func (s *AuthServices) sign(payload string) string {
	mac := hmac.New(sha256.New, s.secretKey)
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func (s *AuthServices) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("Kullanıcı bulunamadı")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("email veya şifre hatalı")
	}

	payload := tokenPayload{
		UserID: user.ID,
		Role:   user.Role,
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadBytes)

	signature := s.sign(encodedPayload)
	token := encodedPayload + "." + signature
	return token, nil
}

func (s *AuthServices) VerifyToken(ctx context.Context, tokenString string) (int, string, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 2 {
		return 0, "", errors.New("gecersiz token formati")
	}
	encodedPayload, signature := parts[0], parts[1]
	expectedSignature := s.sign(encodedPayload)
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return 0, "", errors.New("gecersiz imza")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return 0, "", errors.New("payload decode edilmedi")
	}

	var payload tokenPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return 0, "", errors.New("payload parse edilmedi")
	}

	if time.Now().Unix() > payload.Exp {
		return 0, "", errors.New("token suresi dolmus")
	}

	blaclisted, err := s.authRepo.IsTokenBlackList(ctx, signature)
	if err != nil {
		return 0, "", errors.New("token dogrulanmadi")
	}
	if blaclisted {
		return 0, "", errors.New("token gecersiz kilinmis")
	}

	return payload.UserID, payload.Role, nil

}

func (s *AuthServices) Logout(ctx context.Context, tokenString string) error {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 2 {
		return errors.New("gecersiz token formatı")

	}
	encodedPayload, signature := parts[0], parts[1]

	expectedSignature := s.sign(encodedPayload)
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return errors.New("gecersiz token")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return errors.New("payload decode edilemedi")
	}
	var payload tokenPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return errors.New("payload parse edilemedi")
	}

	remaining := time.Until(time.Unix(payload.Exp, 0))
	if remaining <= 0 {
		return nil
	}
	return s.authRepo.BlackListToken(ctx, signature, remaining)
}
