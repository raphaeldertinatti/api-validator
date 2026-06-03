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

type DiferimentoRepository interface {
	FindByCode(ctx context.Context, ncmCode string) ([]domains.DiferimentoDocument, error)
}

type DiferimentoService struct {
	ncmRepo NCMRepository
	repo    DiferimentoRepository
	gemini  *GeminiService
}

func NewDiferimentoService(ncmRepo NCMRepository, repo DiferimentoRepository, gemini *GeminiService) *DiferimentoService {
	return &DiferimentoService{
		ncmRepo: ncmRepo,
		repo:    repo,
		gemini:  gemini,
	}
}

func (s *DiferimentoService) ValidateDiferimento(ctx context.Context, req domains.ValidateRequest) (*domains.DiferimentoValidacaoResponse, error) {
	// 1. Buscar possíveis diferimentos para a NCM no MongoDB
	docs, err := s.repo.FindByCode(ctx, req.NCM)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar diferimentos: %w", err)
	}

	ncmDescricao, err := s.ncmRepo.FindDescricaoCompleta(ctx, req.NCM)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar NCM: %w", err)
	}

	if len(docs) == 0 {
		return &domains.DiferimentoValidacaoResponse{
			Diferido:      false,
			Justificativa: "Produto não possui diferimento previsto na legislação.",
		}, nil
	}

	// 2. Analisar se a descrição do produto bate com a descrição de algum diferimento
	return s.AnaliseDiferimento(ctx, req, docs, ncmDescricao)
}

func (s *DiferimentoService) AnaliseDiferimento(ctx context.Context, req domains.ValidateRequest, docs []domains.DiferimentoDocument, ncmDescricao string) (*domains.DiferimentoValidacaoResponse, error) {
	type DiferimentoBrief struct {
		ID             interface{}                  `json:"id"`
		Artigo         string                       `json:"artigo"`
		Inciso         *string                      `json:"inciso"`
		Titulo         string                       `json:"titulo"`
		DescricaoLegal string                       `json:"descricao_legal"`
		Produtos       []string                     `json:"produtos"`
		Condicoes      domains.DiferimentoCondicoes `json:"condicoes"`
	}
	var possibleDiferimentos []DiferimentoBrief
	for _, d := range docs {
		possibleDiferimentos = append(possibleDiferimentos, DiferimentoBrief{
			ID:             d.ID,
			Artigo:         d.Artigo,
			Inciso:         d.Inciso,
			Titulo:         d.Titulo,
			DescricaoLegal: d.DescricaoLegal,
			Produtos:       d.Produtos,
			Condicoes:      d.Condicoes,
		})
	}
	diferimentosJSON, _ := json.Marshal(possibleDiferimentos)

	prompt := fmt.Sprintf(`Como especialista tributário, analise se o produto abaixo se enquadra em alguma das previsões de DIFERIMENTO DE ICMS fornecidas.
Considere a descrição do produto e as condições/descrições legais.

Produto: %s
NCM: %s
Descrição NCM: %s
Diferimentos Candidatos: %s

Responda APENAS JSON:
{
  "diferido": true|false,
  "artigo": "...",
  "justificativa": "Sua justificativa aqui em no máximo 20 palavras"
  "observacao": "Somente no caso do artigo 391 de pescados retornar: Nas saídas de estabelecimento com CNAE principal 1020-1/01 ou 1020-1/02 — não é diferido"
}

O campo "diferido" deve ser um booleano.`, req.Descricao, req.NCM, ncmDescricao, string(diferimentosJSON))

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

	candidate := gemResp.Candidates[0]
	rawText := strings.TrimSpace(candidate.Content.Parts[0].Text)
	rawText = strings.TrimPrefix(rawText, "```json")
	rawText = strings.TrimPrefix(rawText, "```")
	rawText = strings.TrimSuffix(rawText, "```")
	rawText = strings.TrimSpace(rawText)

	if rawText == "" {
		return nil, fmt.Errorf("Gemini retornou um texto vazio para a análise de diferimento (finishReason: %s)", candidate.FinishReason)
	}

	var llmResult domains.DiferimentoValidacaoResponse
	if err := json.Unmarshal([]byte(rawText), &llmResult); err != nil {
		return nil, fmt.Errorf("erro ao processar decisão da IA: %w (finishReason: %s, raw: %s)", err, candidate.FinishReason, rawText)
	}

	return &llmResult, nil
}
