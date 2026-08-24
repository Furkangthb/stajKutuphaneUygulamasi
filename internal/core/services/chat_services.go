package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
)

type Chatservices struct {
	apikey       string
	httpsClient  *http.Client
	bookServices *BookServices
}

func NewChatServices(apikey string, bookServices *BookServices) *Chatservices {
	return &Chatservices{apikey: apikey, httpsClient: &http.Client{Timeout: 60 * time.Second}, bookServices: bookServices}
}

const sistemTalimati = `Sen bir kütüphane uygulamasının yardımcı botusun.Sorulara türkçe cevap verirsin.Resmi bir dilde konuşursun.Bu uygulama ile alakasız konuları,bu konu hakkımda bilgim yok diyerek cevap verirsin.Asla iç düşünce sürecini, analizlerini veya adımlarını (Chain of Thought) çıktı olarak verme. İngilizce açıklama yapma. Sadece ve sadece kullanıcıya vereceğin nihai Türkçe cevabı tek bir metin halinde üret.

ÖNEMLİ: Kitaplarla ilgili sorularda SADECE sana mesajın içinde verilen kütüphane kataloğu listesindeki bilgileri kullan. Kendi genel bilgini veya eğitim verindeki bilgiyi ASLA kullanma. Kullanıcının sorduğu kitap sana verilen listede yoksa, kesinlikle "Bu kitap şu an kütüphanemizin sisteminde kayıtlı görünmüyor" de. Bir kitabın var olduğunu, stokta olduğunu veya rezerve edilebilir olduğunu, o kitap sana verilen listede açıkça yer almıyorsa ASLA söyleme. Listede olmayan bir kitap hakkında tahminde bulunma veya varsayımda bulunma.`

// const geminiEndpoint = "https://generativelanguage.googleapis.com/v1beta/models/gemma-4-31b-it:generateContent"
const geminiEndpoint = "https://generativelanguage.googleapis.com/v1beta/models/gemini-3.1-flash-lite:generateContent"

type geminiIstek struct {
	SystemInstruction geminiIcerik   `json:"system_instruction"`
	Contents          []geminiIcerik `json:"contents"`
}

type geminiIcerik struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiYanit struct {
	Candidates []struct {
		Content geminiIcerik `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (s Chatservices) Sohbet(ctx context.Context, kullaniciMesaji string) (string, error) {
	if s.apikey == "" {
		return "", errors.New("api anahtari bos")
	}
	if kullaniciMesaji == "" {
		return "", errors.New("kullanici bos olamaz")
	}

	keywords := kelimeCikar(kullaniciMesaji)
	ilgiliKitaplar, err := s.bookServices.BookSearch(ctx, keywords, 10)
	aramaBulamadi := err != nil || len(ilgiliKitaplar) == 0

	if aramaBulamadi {
		ilgiliKitaplar, err = s.bookServices.BookList(ctx, 20, 0)
		if err != nil {
			ilgiliKitaplar = nil
		}
	}

	prompMetni := kullaniciMesaji
	if len(ilgiliKitaplar) > 0 {
		if aramaBulamadi {
			prompMetni = fmt.Sprintf("Kullanıcının aradığı kitapla eşleşen bir sonuç bulunamadı. Kütüphanenin genel kataloğundan bir örnek liste:\n%s\n\nKullanıcının sorusu: %s\n\nBu liste kütüphanenin TÜM kataloğu değildir, sadece bir örnektir. Kullanıcının sorduğu kitabın bu listede olmaması, kütüphanede olmadığı anlamına gelmez ama kesinlikle 'stokta var' da diyemezsin - bunun yerine 'aradığınız kitabı sistemde bulamadım, lütfen tam adını kontrol edin veya kütüphane personeline danışın' gibi bir cevap ver.",
				kitaplariMetneCevir(ilgiliKitaplar), kullaniciMesaji)
		} else {
			prompMetni = fmt.Sprintf("Kütüphanenin GÜNCEL kataloğu (bu listenin dışında hiçbir kitap kütüphanede mevcut değildir):\n%s\n\nKullanıcının sorusu: %s",
				kitaplariMetneCevir(ilgiliKitaplar), kullaniciMesaji)
		}
	}

	istek := geminiIstek{
		SystemInstruction: geminiIcerik{
			Parts: []geminiPart{{Text: sistemTalimati}},
		},
		Contents: []geminiIcerik{
			{
				Role:  "user",
				Parts: []geminiPart{{Text: prompMetni}},
			},
		},
	}
	body, err := json.Marshal(istek)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s?key=%s", geminiEndpoint, s.apikey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.httpsClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var yanit geminiYanit

	if err := json.Unmarshal(respBody, &yanit); err != nil {
		return "", fmt.Errorf("yanit parse edilemedi")
	}
	if yanit.Error != nil {

		return "", fmt.Errorf("gemini api hatasi: %s", yanit.Error.Message)

	}
	if len(yanit.Candidates) == 0 || len(yanit.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("geminiden bos yanit geldi")
	}
	return yanit.Candidates[0].Content.Parts[0].Text, nil
}

func kelimeCikar(mesaj string) []string {
	durakKelimeler := map[string]bool{
		"bir": true, "bu": true, "şu": true, "ve": true, "ile": true,
		"mi": true, "mı": true, "misin": true, "mısın": true,
		"var": true, "yok": true, "nasıl": true, "hangi": true,
		"kitap": true, "kitabı": true, "kitaplar": true,
	}

	kelimeler := strings.Fields(strings.ToLower(mesaj))
	var sonuc []string
	for _, k := range kelimeler {
		k = strings.Trim(k, ".,?!'\"")
		if len(k) < 3 || durakKelimeler[k] {
			continue
		}
		sonuc = append(sonuc, k)
	}
	return sonuc
}

func kitaplariMetneCevir(kitaplar []*domain.Book) string {
	var satirlar []string
	for _, k := range kitaplar {
		satirlar = append(satirlar, fmt.Sprintf("- %s / %s / Tur: %s / Stok: %d",
			k.Title, k.Author, k.Genre, k.StockCount))
	}
	return strings.Join(satirlar, "\n")
}
