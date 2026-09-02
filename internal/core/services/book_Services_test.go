package services

import (
	"context"
	"testing"
	"time"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
)

type fakeBookRepo struct {
	book *domain.Book

	createErr error
	updateErr error
	deleteErr error
	getErr    error
	listErr   error
	searchErr error

	createdBook *domain.Book
	updatedBook *domain.Book

	listLimit     int
	listOffset    int
	searchLimit   int
	searchWords   []string
	searchResault []*domain.Book
}

func (f *fakeBookRepo) BookCreate(ctx context.Context, b *domain.Book) error {
	f.createdBook = b
	return f.createErr
}

func (f *fakeBookRepo) BookUpdate(ctx context.Context, b *domain.Book) error {
	f.updatedBook = b
	return f.updateErr

}

func (f *fakeBookRepo) BookDelete(ctx context.Context, id int64) error {
	return f.deleteErr
}

func (f *fakeBookRepo) BookGetByID(ctx context.Context, id int64) (*domain.Book, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.book, nil
}

func (f *fakeBookRepo) BookList(ctx context.Context, limit, offset int) ([]*domain.Book, error) {
	f.listLimit = limit
	f.listOffset = offset
	if f.listErr != nil {
		return nil, f.listErr
	}
	return []*domain.Book{f.book}, nil
}

func (f *fakeBookRepo) BookSearch(ctx context.Context, limit int, keywords []string) ([]*domain.Book, error) {
	f.searchLimit = limit
	f.searchWords = keywords

	if f.searchErr != nil {
		return nil, f.searchErr
	}
	return f.searchResault, nil
}

func TestBookAdd_Basarali(t *testing.T) {
	fakeRepo := &fakeBookRepo{}
	s := NewBookServices(fakeRepo)
	yayinTarihi := time.Now()
	book, err := s.BookAdd(context.Background(), "5555", "deneme", "yazar", "aksiyon", yayinTarihi, "bu bir denemedir")
	if err != nil {
		t.Fatalf("Beklenmeyen hata %v", err)
	}
	if book.ISBN != "5555" || book.Author != "yazar" || book.Title != "deneme" {
		t.Errorf("alanlar eslesmiyo %+v", book)
	}

	if fakeRepo.createdBook == nil {
		t.Fatalf("repoya erisemiyor")
	}
	if fakeRepo.createdBook.ISBN != "5555" {
		t.Errorf("ISBN degerleri eslemiyor ,beklenen=5555 gelen=%v", fakeRepo.createdBook.ISBN)
	}
}

func TestBookAdd_Hata(t *testing.T) {
	fakeRepo := &fakeBookRepo{createErr: errTest("hata")}
	s := NewBookServices(fakeRepo)
	_, err := s.BookAdd(context.Background(), "5555", "deneme", "yazar", "aksiyon", time.Now(), "bu bir denemedir")
	if err == nil {
		t.Fatalf("Hata bekleniyordu,nil dondu")
	}
}

func TestBookGet_Basarali(t *testing.T) {
	beklenenKitap := &domain.Book{ID: 6, Title: "deneme"}
	fakeRepo := &fakeBookRepo{book: beklenenKitap}
	s := NewBookServices(fakeRepo)
	book, err := s.BookGet(context.Background(), 6)
	if err != nil {
		t.Fatalf("Beklenmeyen hata %v", err)
	}
	if book.ID != 6 {
		t.Errorf("id ler eslemiyor, beklenen=6  gelen=%d", book.ID)
	}
}

func TestBookGet_Hata(t *testing.T) {
	fakeRepo := &fakeBookRepo{getErr: errTest("hata")}
	s := NewBookServices(fakeRepo)
	_, err := s.BookGet(context.Background(), 6)
	if err == nil {
		t.Fatalf("Hata bekleniyordu,nil dondu")
	}
}

func TestBookDelete_Basarali(t *testing.T) {
	fakeRepo := &fakeBookRepo{}
	s := NewBookServices(fakeRepo)

	err := s.BookDelete(context.Background(), 6)
	if err != nil {
		t.Fatalf("Beklenmeyen hata")
	}

}

func TestBookDelete_Hata(t *testing.T) {
	fakeRepo := &fakeBookRepo{deleteErr: errTest("hata")}
	s := NewBookServices(fakeRepo)
	err := s.BookDelete(context.Background(), 2)
	if err == nil {
		t.Fatalf("Hata bekleniyordu,nil dondu")
	}
}

func TestBookUpdate_Basarali(t *testing.T) {
	fakeRepo := &fakeBookRepo{}
	s := NewBookServices(fakeRepo)

	book, err := s.BookUpdate(context.Background(), 5, "1111", "deneme", "yazar", "fantastik", time.Now(), "denemedir.")
	if err != nil {
		t.Fatalf("Beklenmeyen hata olustu %v", err)
	}
	if book.ID != 5 || book.Title != "deneme" {
		t.Errorf("Guncellenen alanlar yanlis")
	}
	if fakeRepo.updatedBook == nil {
		t.Fatalf("repo.update hic cagırılmadi")
	}
	if fakeRepo.updatedBook.ID != 5 {
		t.Fatalf("id ler eslesmiyor , beklenen=5 gelen=%v", fakeRepo.updatedBook.ID)
	}
}

func TestBookUpdate_Hata(t *testing.T) {
	fakeRepo := &fakeBookRepo{updateErr: errTest("hata")}
	s := NewBookServices(fakeRepo)
	_, err := s.BookUpdate(context.Background(), 5, "5555", "deneme", "yazar", "aksiyon", time.Now(), "denemdir")
	if err == nil {
		t.Fatalf("hata bekleniyordu,nil dondu")
	}
}

func TestBookList_Basarali(t *testing.T) {
	fakeRepo := &fakeBookRepo{}
	s := NewBookServices(fakeRepo)
	_, err := s.BookList(context.Background(), 1, 20)
	if err != nil {
		t.Fatalf("Beklenmeyen bir hata olustu  %v", err)
	}

}
func FuzzBookList_Sayfalama(f *testing.F) {
	f.Add(1, 10)
	f.Add(0, 5)

	f.Fuzz(func(t *testing.T, page int, pageSize int) {

		fakeRepo := &fakeBookRepo{}
		s := NewBookServices(fakeRepo)

		_, _ = s.BookList(context.Background(), page, pageSize)

		if fakeRepo.listOffset < 0 {
			t.Errorf("Fuzzer bir açık buldu! Negatif offset oluştu. Page: %v, Size: %v, Olusan Offset: %v",
				page, pageSize, fakeRepo.listOffset)
		}
	})
}
func TestBookList_Hata(t *testing.T) {
	fakeRepo := &fakeBookRepo{listErr: errTest("hata")}
	s := NewBookServices(fakeRepo)
	_, err := s.BookList(context.Background(), 5, 20)
	if err == nil {
		t.Fatalf("hata beklenirken nil dondu")
	}
}

func TestBookSearch(t *testing.T) {
	beklenenSonuc := []*domain.Book{{ID: 5, Title: "Deneme"}}
	fakeRepo := &fakeBookRepo{searchResault: beklenenSonuc}
	s := NewBookServices(fakeRepo)
	books, err := s.BookSearch(context.Background(), []string{"aksiyon", "fantastik"}, 10)
	if err != nil {
		t.Fatalf("Beklenmeyen bir hata olustu")
	}
	if books[0].Title != "Deneme" {
		t.Errorf("beklenmeyen kitap")
	}
	if fakeRepo.searchLimit != 10 {
		t.Errorf("limitler eslemiyor, beklenen=10 gelen=%v", fakeRepo.searchLimit)
	}
}

func TestBookSearch_Hata(t *testing.T) {
	fakeRepo := &fakeBookRepo{searchErr: errTest("hata")}
	s := NewBookServices(fakeRepo)
	_, err := s.BookSearch(context.Background(), []string{"aksiyon"}, 10)
	if err == nil {
		t.Fatalf("Hata bekleniyordu nil dondu")
	}

}
