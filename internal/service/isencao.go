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

type IsencaoRepository interface {
	FindByCode(ctx context.Context, ncmCode string) ([]domains.IsencaoDocument, error)
}

type IsencaoService struct {
	ncmRepo NCMRepository
	repo    IsencaoRepository
	gemini  *GeminiService
}

func NewIsencaoService(ncmRepo NCMRepository, repo IsencaoRepository, gemini *GeminiService) *IsencaoService {
	return &IsencaoService{
		ncmRepo: ncmRepo,
		repo:    repo,
		gemini:  gemini,
	}
}

func (s *IsencaoService) ValidateIsencao(ctx context.Context, req domains.ValidateRequest) (*domains.IsencaoValidacaoResponse, error) {
	// 1. Buscar possíveis isenções para a NCM no MongoDB
	docs, err := s.repo.FindByCode(ctx, req.NCM)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar isenções: %w", err)
	}

	ncmDescricao, err := s.ncmRepo.FindDescricaoCompleta(ctx, req.NCM)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar NCM: %w", err)
	}

	if len(docs) == 0 {
		return &domains.IsencaoValidacaoResponse{
			Isento:        false,
			Justificativa: "Produto não possui isenções previstas na legislação.",
		}, nil
	}

	// 2. Analisar se a descrição do produto bate com a descrição de alguma isenção
	return s.AnaliseIsencao(ctx, req, docs, ncmDescricao)
}

func (s *IsencaoService) AnaliseIsencao(ctx context.Context, req domains.ValidateRequest, docs []domains.IsencaoDocument, ncmDescricao string) (*domains.IsencaoValidacaoResponse, error) {
	type IsencaoBrief struct {
		ID             interface{}              `json:"id"`
		Artigo         string                   `json:"artigo"`
		Paragrafo      string                   `json:"paragrafo"`
		Inciso         string                   `json:"inciso"`
		Titulo         string                   `json:"titulo"`
		DescricaoLegal string                   `json:"descricao_legal"`
		Produtos       []string                 `json:"produtos"`
		Condicoes      domains.IsencaoCondicoes `json:"condicoes"`
	}
	var possibleIsencoes []IsencaoBrief
	for _, d := range docs {
		possibleIsencoes = append(possibleIsencoes, IsencaoBrief{
			ID:             d.ID,
			Artigo:         d.Artigo,
			Paragrafo:      d.Paragrafo,
			Inciso:         d.Inciso,
			Titulo:         d.Titulo,
			DescricaoLegal: d.DescricaoLegal,
			Produtos:       d.Produtos,
			Condicoes:      d.Condicoes,
		})
	}
	isencoesJSON, _ := json.Marshal(possibleIsencoes)

	// Prompt focado na comparação de descrições e escolha entre candidatos
	prompt := fmt.Sprintf(`Como especialista tributário, analise se o produto abaixo se enquadra em alguma das previsões de isenção fornecidas.
Cada documento candidato é independente e suas condições não se relacionam, escolha a mais específica ou adequada de acordo com a descrição do produto e as previsões legais.

Produto: %s
NCM: %s
Descrição NCM: %s
Isenções Candidatas: %s

Responda APENAS JSON:
{
  "isento": true|false,
  "artigo": "...",
  "paragrafo": "...",
  "justificativa": "Sua justificativa aqui em no máximo 20 palavras"
}

O campo "isento" deve ser um booleano (true ou false) e não uma string.`, req.Descricao, req.NCM, ncmDescricao, string(isencoesJSON))

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
		return nil, fmt.Errorf("Gemini retornou um texto vazio para a análise de isenção")
	}

	var llmResult domains.IsencaoValidacaoResponse
	if err := json.Unmarshal([]byte(rawText), &llmResult); err != nil {
		return nil, fmt.Errorf("erro ao processar decisão da IA: %w (raw: %s)", err, rawText)
	}

	return &llmResult, nil
}
