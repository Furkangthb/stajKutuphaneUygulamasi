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
	chatRepo     domain.ChatRepository
}

func NewChatServices(apikey string, bookServices *BookServices, chatRepo domain.ChatRepository) *Chatservices {
	return &Chatservices{
		apikey:       apikey,
		httpsClient:  &http.Client{Timeout: 60 * time.Second},
		bookServices: bookServices,
		chatRepo:     chatRepo,
	}
}

const sistemTalimati = `Sen bir kütüphane uygulamasının akıllı asistanısın. Türkçe ve resmi bir dilde cevap verirsin.
Eğer kullanıcı kütüphanedeki bir kitabı, özetini, yazarını veya stok durumunu soruyorsa, kesinlikle sana sağlanan 'KitapAra' fonksiyonunu kullanarak veritabanında arama yapmalısın. Geçmiş mesajlarda bahsedilen bir kitabın özetini istese bile bu fonksiyonu o kitabın adıyla tekrar çağır.
Fonksiyondan gelen sonuca göre kullanıcıya bilgi ver. Kullanıcının aradığı kitap kütüphanede yoksa "Sistemimizde bulamadım" de.
Asla uydurma kitap bilgisi verme.Senin rezerve yetkin yok ama kullanıcıya bilgi verirsin. Eğer kullanıcı rezerve etmek isterse, "Rezerve etmek için lütfen kütüphane uygulamasına giriş yapın" de.`

const geminiEndpoint = "https://generativelanguage.googleapis.com/v1beta/models/gemini-3.1-flash-lite:generateContent"

type geminiIstek struct {
	SystemInstruction *geminiIcerik  `json:"system_instruction,omitempty"`
	Contents          []geminiIcerik `json:"contents"`
	Tools             []geminiTool   `json:"tools,omitempty"`
}

type geminiIcerik struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart map[string]interface{}

type geminiTool struct {
	FunctionDeclarations []geminiFunctionDeclaration `json:"functionDeclarations"`
}

type geminiFunctionDeclaration struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Parameters  geminiSchema `json:"parameters"`
}

type geminiSchema struct {
	Type        string                  `json:"type"`
	Properties  map[string]geminiSchema `json:"properties,omitempty"`
	Required    []string                `json:"required,omitempty"`
	Items       *geminiSchema           `json:"items,omitempty"`
	Description string                  `json:"description,omitempty"`
}

type geminiYanit struct {
	Candidates []struct {
		Content geminiIcerik `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (s Chatservices) Sohbet(ctx context.Context, userID int, kullaniciMesaji string) (string, error) {
	if s.apikey == "" {
		return "", errors.New("api anahtari bos")
	}
	if kullaniciMesaji == "" {
		return "", errors.New("kullanici bos olamaz")
	}

	var tumIcerikler []geminiIcerik

	if s.chatRepo != nil {
		gecmis, err := s.chatRepo.GetHistory(ctx, userID, 10)
		if err == nil {
			for i := len(gecmis) - 1; i >= 0; i-- {
				rol := "user"
				if gecmis[i].Role == "assistant" {
					rol = "model"
				}
				tumIcerikler = append(tumIcerikler, geminiIcerik{
					Role:  rol,
					Parts: []geminiPart{{"text": gecmis[i].Message}},
				})
			}
		}
	}

	tumIcerikler = append(tumIcerikler, geminiIcerik{
		Role:  "user",
		Parts: []geminiPart{{"text": kullaniciMesaji}},
	})

	araclar := []geminiTool{
		{
			FunctionDeclarations: []geminiFunctionDeclaration{
				{
					Name:        "KitapAra",
					Description: "Kütüphane veritabanında kitap araması yapar. Kullanıcının aradığı kitabın adını, yazarını veya temasını anahtar kelime olarak alır.",
					Parameters: geminiSchema{
						Type: "OBJECT",
						Properties: map[string]geminiSchema{
							"kelimeler": {
								Type:        "ARRAY",
								Description: "Aranacak kelimeler listesi. Örn: ['Geleceğin', 'Fiziği'] veya ['aşk', 'romanı']",
								Items:       &geminiSchema{Type: "STRING"},
							},
						},
						Required: []string{"kelimeler"},
					},
				},
			},
		},
	}

	istek := geminiIstek{
		SystemInstruction: &geminiIcerik{Parts: []geminiPart{{"text": sistemTalimati}}},
		Contents:          tumIcerikler,
		Tools:             araclar,
	}

	yanit, err := s.geminiIstekAt(ctx, istek)
	if err != nil {
		return "", err
	}

	modelContent := yanit.Candidates[0].Content
	var functionCallPart map[string]interface{}
	var nihaiCevapMetni string

	for _, part := range modelContent.Parts {
		if fc, ok := part["functionCall"].(map[string]interface{}); ok {
			functionCallPart = fc
			break
		}
	}

	if functionCallPart != nil {
		fonksiyonAdi, _ := functionCallPart["name"].(string)

		if fonksiyonAdi == "KitapAra" {
			var aranacakKelimeler []string
			if args, ok := functionCallPart["args"].(map[string]interface{}); ok {
				if kelimelerInterface, ok2 := args["kelimeler"].([]interface{}); ok2 {
					for _, k := range kelimelerInterface {
						aranacakKelimeler = append(aranacakKelimeler, fmt.Sprintf("%v", k))
					}
				}
			}

			fmt.Printf("Yapay Zeka Arıyor: %v\n", aranacakKelimeler)
			ilgiliKitaplar, _ := s.bookServices.BookSearch(ctx, aranacakKelimeler, 10)

			aramaSonucuMetni := "Aranan kriterlere uygun kitap bulunamadı."
			if len(ilgiliKitaplar) > 0 {
				aramaSonucuMetni = kitaplariMetneCevir(ilgiliKitaplar)
			}

			tumIcerikler = append(tumIcerikler, modelContent)

			tumIcerikler = append(tumIcerikler, geminiIcerik{
				Role: "user",
				Parts: []geminiPart{
					{
						"functionResponse": map[string]interface{}{
							"name": "KitapAra",
							"response": map[string]interface{}{
								"sonuc": aramaSonucuMetni,
							},
						},
					},
				},
			})

			ikinciIstek := geminiIstek{
				SystemInstruction: &geminiIcerik{Parts: []geminiPart{{"text": sistemTalimati}}},
				Contents:          tumIcerikler,
			}

			ikinciYanit, err := s.geminiIstekAt(ctx, ikinciIstek)
			if err != nil {
				return "", err
			}

			for _, part := range ikinciYanit.Candidates[0].Content.Parts {
				if txt, ok := part["text"].(string); ok {
					nihaiCevapMetni += txt
				}
			}
		}
	} else {
		for _, part := range modelContent.Parts {
			if txt, ok := part["text"].(string); ok {
				nihaiCevapMetni += txt
			}
		}
	}

	if strings.TrimSpace(nihaiCevapMetni) == "" {
		nihaiCevapMetni = "Cevap üretilirken bir hata oluştu, lütfen tekrar deneyin."
	}

	if s.chatRepo != nil {
		s.chatRepo.Save(ctx, domain.ChatMessage{UserID: userID, Role: "user", Message: kullaniciMesaji})
		s.chatRepo.Save(ctx, domain.ChatMessage{UserID: userID, Role: "assistant", Message: nihaiCevapMetni})
	}

	return nihaiCevapMetni, nil
}

func (s Chatservices) geminiIstekAt(ctx context.Context, istek geminiIstek) (*geminiYanit, error) {
	body, err := json.Marshal(istek)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s?key=%s", geminiEndpoint, s.apikey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpsClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var yanit geminiYanit
	if err := json.Unmarshal(respBody, &yanit); err != nil {
		return nil, fmt.Errorf("yanit parse edilemedi: %s", string(respBody))
	}
	if yanit.Error != nil {
		return nil, fmt.Errorf("gemini api hatasi: %s", yanit.Error.Message)
	}
	if len(yanit.Candidates) == 0 || len(yanit.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("geminiden bos yanit geldi")
	}

	return &yanit, nil
}

func kitaplariMetneCevir(kitaplar []*domain.Book) string {
	var satirlar []string
	for _, k := range kitaplar {
		durum := "Müsait değil / rezerve"
		if k.Available {
			durum = "Müsait"
		}
		ozet := k.Description
		if ozet == "" {
			ozet = "(özet yok)"
		}
		satirlar = append(satirlar, fmt.Sprintf("- %s / %s / Tur: %s / Durum: %s / Ozet: %s",
			k.Title, k.Author, k.Genre, durum, ozet))
	}
	return strings.Join(satirlar, "\n")
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

func (s Chatservices) GetHistory(ctx context.Context, userID int, limit int) ([]domain.ChatMessage, error) {
	if s.chatRepo == nil {
		return nil, errors.New("chat gecmisi kullanilamiyor")
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	gecmis, err := s.chatRepo.GetHistory(ctx, userID, limit)
	if err != nil {
		return nil, err
	}
	for i, j := 0, len(gecmis)-1; i < j; i, j = i+1, j-1 {
		gecmis[i], gecmis[j] = gecmis[j], gecmis[i]
	}
	return gecmis, nil
}
