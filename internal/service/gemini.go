package service

import (
	"net/http"
	"time"
)

// GeminiService encapsula a chamada à API do Gemini
type GeminiService struct {
	apiKey     string
	httpClient *http.Client
	model      string
}

func NewGeminiService(apiKey string) *GeminiService {
	return &GeminiService{
		apiKey: apiKey,
		model:  "gemini-3-flash-preview", // melhor custo-benefício no plano gratuito
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// structs para a API REST do Gemini
type geminiRequest struct {
	Contents         []geminiContent        `json:"contents"`
	GenerationConfig geminiGenerationConfig `json:"generationConfig"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiGenerationConfig struct {
	Temperature     float64 `json:"temperature"`
	MaxOutputTokens int     `json:"maxOutputTokens"`
}

type geminiResponse struct {
	Candidates []struct {
		Content geminiContent `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error,omitempty"`
}
