package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Chatservices struct {
	apikey      string
	httpsClient *http.Client
}

func NewChatServices(apikey string) *Chatservices {
	return &Chatservices{apikey: apikey, httpsClient: &http.Client{Timeout: 20 * time.Second}}
}

const sistemTalimati = `Sen bir kütüphane uygulamasının yardımcı botusun.Sorulara türkçe cevap verirsin.Resmi bir dilde konuşursun.Bu uygulama ile alakasız konuları,bu konu hakkımda bilgim yok diyerek cevap verirsin.`

const geminiEndpoint = "https://generativelanguage.googleapis.com/v1beta/models/gemini-3.6-flash:generateContent"

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
	istek := geminiIstek{
		SystemInstruction: geminiIcerik{
			Parts: []geminiPart{{Text: sistemTalimati}},
		},
		Contents: []geminiIcerik{
			{
				Role:  "user",
				Parts: []geminiPart{{Text: kullaniciMesaji}},
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
