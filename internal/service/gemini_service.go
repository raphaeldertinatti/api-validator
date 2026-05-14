package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// GeminiClient encapsula a chamada à API do Gemini
type GeminiClient struct {
	apiKey     string
	httpClient *http.Client
	model      string
}

func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{
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

// ValidateNCM verifica se a descrição do produto condiz com a NCM
func (g *GeminiClient) ValidateNCM(ctx context.Context, prodDesc, ncmDesc string) (string, bool, error) {
	prompt := fmt.Sprintf(`Você é um especialista em classificação fiscal (NCM).
Analise se a descrição do produto fornecida condiz com a descrição oficial da NCM.

Descrição do Produto: %s
Descrição Oficial da NCM: %s

Responda EXATAMENTE neste formato JSON:
{"status": "SIM" ou "NAO", "justificativa": "motivo em até 15 palavras"}

Responda APENAS o JSON, sem markdown, sem código, sem explicações adicionais.`, prodDesc, ncmDesc)

	reqBody := geminiRequest{
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: prompt}}},
		},
		GenerationConfig: geminiGenerationConfig{
			Temperature:     0.1, // mais determinístico para validação
			MaxOutputTokens: 1024,
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", false, fmt.Errorf("erro ao serializar request: %w", err)
	}

	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		g.model, g.apiKey,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return "", false, fmt.Errorf("erro ao criar request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", false, fmt.Errorf("erro na chamada HTTP: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	var gemResp geminiResponse
	if err := json.Unmarshal(respBody, &gemResp); err != nil {
		return "", false, fmt.Errorf("erro ao deserializar resposta: %w", err)
	}

	// Verifica erro da API
	if gemResp.Error != nil {
		return "", false, fmt.Errorf("erro da API Gemini [%d]: %s", gemResp.Error.Code, gemResp.Error.Message)
	}

	if len(gemResp.Candidates) == 0 || len(gemResp.Candidates[0].Content.Parts) == 0 {
		return "", false, fmt.Errorf("resposta vazia do Gemini")
	}

	rawText := gemResp.Candidates[0].Content.Parts[0].Text
	rawText = strings.TrimSpace(rawText)

	// Parse do JSON de resposta
	var result struct {
		Status        string `json:"status"`
		Justificativa string `json:"justificativa"`
	}
	if err := json.Unmarshal([]byte(rawText), &result); err != nil {
		// Fallback: parse simples se o modelo não retornar JSON perfeito
		isValid := strings.Contains(strings.ToUpper(rawText), `"SIM"`)
		return rawText, isValid, nil
	}

	respFormatada := fmt.Sprintf("STATUS: %s\nJUSTIFICATIVA: %s", result.Status, result.Justificativa)
	isValid := strings.EqualFold(result.Status, "SIM")

	return respFormatada, isValid, nil
}
