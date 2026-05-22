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

type CESTRepository interface {
	FindByCode(ctx context.Context, ncmCode string) ([]domains.CESTDocument, error)
}

type CESTService struct {
	repo   CESTRepository
	gemini *GeminiService
}

func NewCESTService(repo CESTRepository, gemini *GeminiService) *CESTService {
	return &CESTService{
		repo:   repo,
		gemini: gemini,
	}
}

func (s *CESTService) ValidateCEST(ctx context.Context, req domains.ValidateRequest) (*domains.CESTValidacaoResponse, error) {
	// 1. Buscar possíveis CESTs para a NCM no MongoDB
	docs, err := s.repo.FindByCode(ctx, req.NCM)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar CEST: %w", err)
	}

	if len(docs) == 0 {
		return &domains.CESTValidacaoResponse{
			Status:        "NAO_ENQUADRADO",
			Justificativa: "NCM não possui códigos CEST previstos na legislação.",
		}, nil
	}

	// 2. Analisar se a descrição do produto bate com a descrição de algum CEST
	return s.AnaliseCEST(ctx, req, docs)
}

func (s *CESTService) AnaliseCEST(ctx context.Context, req domains.ValidateRequest, docs []domains.CESTDocument) (*domains.CESTValidacaoResponse, error) {
	type CestBrief struct {
		CEST      string `json:"cest"`
		Descricao string `json:"descricao"`
		Segmento  string `json:"segmento"`
	}
	var possibleCests []CestBrief
	for _, d := range docs {
		possibleCests = append(possibleCests, CestBrief{
			CEST:      d.CEST,
			Descricao: d.Descricao,
			Segmento:  d.Segmento.Descricao,
		})
	}
	cestsJSON, _ := json.Marshal(possibleCests)

	// Prompt focado na comparação de descrições
	prompt := fmt.Sprintf(`Como especialista tributário, valide se o produto abaixo se enquadra em algum dos códigos CEST da lista. 
Nota: O enquadramento no CEST depende da descrição do produto ser compatível com a descrição específica do CEST, não apenas da NCM.

Produto: %s
CESTs Candidatos: %s

Responda APENAS JSON:
{
  "status": "DEFINIDO"|"AMBIGUO"|"NAO_ENQUADRADO",
  "enquadrado": {"cest":"", "descricao":"", "segmento":""},
  "possibilidades": [],
  "justificativa": "Sua justificativa aqui em no máximo 20 palavras"
}

DEFINIDO: Descrição do produto bate com 1 CEST.
AMBIGUO: Descrição bate com >1 CEST.
NAO_ENQUADRADO: Descrição do produto não bate com as descrições específicas dos CESTs listados.`, req.Descricao, string(cestsJSON))

	reqBody := geminiRequest{
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: prompt}}},
		},
		GenerationConfig: geminiGenerationConfig{
			Temperature:     0.1,
			MaxOutputTokens: 1024,
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

	hReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar request HTTP: %w", err)
	}
	hReq.Header.Set("Content-Type", "application/json")

	respHTTP, err := s.gemini.httpClient.Do(hReq)
	if err != nil {
		return nil, fmt.Errorf("erro na chamada HTTP ao Gemini: %w", err)
	}
	defer respHTTP.Body.Close()

	respBody, err := io.ReadAll(respHTTP.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta do Gemini: %w", err)
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

	rawText := strings.TrimSpace(gemResp.Candidates[0].Content.Parts[0].Text)
	rawText = strings.TrimPrefix(rawText, "```json")
	rawText = strings.TrimPrefix(rawText, "```")
	rawText = strings.TrimSuffix(rawText, "```")
	rawText = strings.TrimSpace(rawText)

	if rawText == "" {
		return nil, fmt.Errorf("Gemini retornou um texto vazio para a análise de CEST")
	}

	var llmResult domains.CESTValidacaoResponse
	if err := json.Unmarshal([]byte(rawText), &llmResult); err != nil {
		return nil, fmt.Errorf("erro ao processar decisão da IA: %w (raw: %s)", err, rawText)
	}

	return &llmResult, nil
}
