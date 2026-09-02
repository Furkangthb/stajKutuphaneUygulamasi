package services

import (
	"context"
	"testing"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
)

type fakeChatRepo struct {
	messages      []domain.ChatMessage
	getHistoryErr error

	capturedUserID int
	capturedLimit  int
}

func (f *fakeChatRepo) Save(ctx context.Context, msq domain.ChatMessage) error {
	f.messages = append(f.messages, msq)
	return nil
}

func (f *fakeChatRepo) GetHistory(ctx context.Context, userID int, limit int) ([]domain.ChatMessage, error) {
	f.capturedUserID = userID
	f.capturedLimit = limit
	if f.getHistoryErr != nil {
		return nil, f.getHistoryErr
	}
	return f.messages, nil
}

func TestChatGetHistory_Basarali(t *testing.T) {
	fakeRepo := &fakeChatRepo{}
	s := NewChatServices("key", nil, fakeRepo)
	_, err := s.GetHistory(context.Background(), 8, 10)
	if err != nil {
		t.Fatalf("Beklenmeyen hata %v", err)
	}

	if fakeRepo.capturedUserID != 8 {
		t.Errorf("id ler eslesmiyor, beklenen=8  gelen=%d", fakeRepo.capturedUserID)
	}
}

func TestChatGetHistory_Hata(t *testing.T) {
	s := NewChatServices("key", nil, nil)
	_,err:=s.GetHistory(context.Background(),5,20)
	if err==nil{
		t.Fatalf("Hata bekleniyordu ,nil dondu")
	}
}
