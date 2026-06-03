package service

import (
	"api-validator/internal/domains"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type NCMRepository interface {
	FindDescricaoCompleta(ctx context.Context, ncmCode string) (string, error)
}

type NCMService struct {
	repo   NCMRepository
	gemini *GeminiService
}

func NewNCMService(repo NCMRepository, gemini *GeminiService) *NCMService {
	return &NCMService{
		repo:   repo,
		gemini: gemini,
	}
}

func (s *NCMService) ValidateNCM(ctx context.Context, req domains.ValidateRequest) (*domains.NCMValidacaoResponse, error) {
	// 1. Buscar descrição do NCM no MongoDB via repositório
	ncmDescricao, err := s.repo.FindDescricaoCompleta(ctx, req.NCM)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar NCM: %w", err)
	}

	if ncmDescricao == "" {
		ncmDescricao = "NCM não encontrado na base de dados."
	}

	// 2. Validar com Gemini usando o novo cliente REST
	compatibilidade, err := s.CompatibilidadeNCM(ctx, req.Descricao, ncmDescricao)
	if err != nil {
		return nil, fmt.Errorf("erro na validação com Gemini: %w", err)
	}

	return &domains.NCMValidacaoResponse{
		NCM:          req.NCM,
		Descricao:    req.Descricao,
		NCMDescricao: ncmDescricao,
		RetornoLLM:   *compatibilidade,
	}, nil
}

// CompatibilidadeNCM verifica se a descrição do produto é compatível com a NCM
func (s *NCMService) CompatibilidadeNCM(ctx context.Context, prodDesc, ncmDesc string) (*domains.NCMCompatibilidadeResult, error) {
	prompt := fmt.Sprintf(`Você é um especialista em classificação fiscal (NCM).
Analise se a descrição do produto fornecida é compatível com a descrição oficial da NCM.

Descrição do Produto: %s
Descrição Oficial da NCM: %s

Responda APENAS este JSON válido, sem markdown, sem explicações adicionais:
{"status": "SIM", "justificativa": "motivo em até 20 palavras"}

O campo status deve ser exatamente "SIM" ou "NAO".`, prodDesc, ncmDesc)

	reqBody := geminiRequest{
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: prompt}}},
		},
		GenerationConfig: geminiGenerationConfig{
			Temperature:      0.1,
			MaxOutputTokens:  2048,
			ResponseMimeType: "application/json",
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar request: %w", err)
	}

	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		s.gemini.model, s.gemini.apiKey,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar request HTTP: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.gemini.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na chamada HTTP: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	var gemResp geminiResponse
	if err := json.Unmarshal(respBody, &gemResp); err != nil {
		return nil, fmt.Errorf("erro ao deserializar resposta Gemini: %w", err)
	}

	if gemResp.Error != nil {
		return nil, fmt.Errorf("erro da API Gemini [%d]: %s", gemResp.Error.Code, gemResp.Error.Message)
	}

	if len(gemResp.Candidates) == 0 || len(gemResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("resposta vazia do Gemini")
	}

	candidate := gemResp.Candidates[0]
	rawText := strings.TrimSpace(candidate.Content.Parts[0].Text)

	// limpa markdown caso o modelo desobedeça o prompt
	rawText = strings.TrimPrefix(rawText, "```json")
	rawText = strings.TrimPrefix(rawText, "```")
	rawText = strings.TrimSuffix(rawText, "```")
	rawText = strings.TrimSpace(rawText)

	if rawText == "" {
		return nil, fmt.Errorf("Gemini retornou um texto vazio para a análise de NCM (finishReason: %s)", candidate.FinishReason)
	}

	var llmResult struct {
		Status        string `json:"status"`
		Justificativa string `json:"justificativa"`
	}
	if err := json.Unmarshal([]byte(rawText), &llmResult); err != nil {
		return nil, fmt.Errorf("erro ao processar decisão da IA: %w (finishReason: %s, raw: %s)", err, candidate.FinishReason, rawText)
	}

	status := strings.ToUpper(strings.TrimSpace(llmResult.Status))
	return &domains.NCMCompatibilidadeResult{
		Compativel:    status == "SIM",
		Status:        status,
		Justificativa: llmResult.Justificativa,
	}, nil
}
