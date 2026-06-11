package service

import (
	"api-validator/internal/domains"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
)

type ReducaoRepository interface {
	FindByCode(ctx context.Context, ncmCode string) ([]domains.ReducaoDocument, error)
}

type ReducaoService struct {
	ncmRepo NCMRepository
	repo    ReducaoRepository
	gemini  *GeminiService
}

func NewReducaoService(ncmRepo NCMRepository, repo ReducaoRepository, gemini *GeminiService) *ReducaoService {
	return &ReducaoService{
		ncmRepo: ncmRepo,
		repo:    repo,
		gemini:  gemini,
	}
}

func (s *ReducaoService) ValidateReducao(aliquotaResult domains.AliquotaValidacaoResponse, ctx context.Context, req domains.ValidateRequest) (*domains.ReducaoValidacaoResponse, error) {
	// 1. Buscar possíveis reduções para a NCM no MongoDB
	docs, err := s.repo.FindByCode(ctx, req.NCM)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar isenções: %w", err)
	}

	ncmDescricao, err := s.ncmRepo.FindDescricaoCompleta(ctx, req.NCM)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar NCM: %w", err)
	}

	if len(docs) == 0 {
		return &domains.ReducaoValidacaoResponse{
			ReducaoBase:   0,
			Justificativa: "Produto não possui reduções previstas na legislação.",
		}, nil
	}

	// 2. Analisar se a descrição do produto bate com a descrição de alguma redução
	reducaoaResult, err := s.AnaliseReducao(ctx, req, docs, ncmDescricao)

	finalResult, err := s.CalcularReducao(ctx, req, reducaoaResult, aliquotaResult)
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular redução: %w", err)
	}

	return finalResult, nil
}

func (s *ReducaoService) AnaliseReducao(ctx context.Context, req domains.ValidateRequest, docs []domains.ReducaoDocument, ncmDescricao string) (*domains.ReducaoValidacaoResponse, error) {
	type ReducaoBrief struct {
		ID             interface{}              `json:"id"`
		Artigo         string                   `json:"artigo"`
		Inciso         string                   `json:"inciso"`
		Titulo         string                   `json:"titulo"`
		DescricaoLegal string                   `json:"descricao_legal"`
		Produtos       []string                 `json:"produtos"`
		Condicoes      domains.ReducaoCondicoes `json:"condicoes"`
	}
	var possibleReducoes []ReducaoBrief
	for _, d := range docs {
		possibleReducoes = append(possibleReducoes, ReducaoBrief{
			ID:             d.ID,
			Artigo:         d.Artigo,
			Inciso:         d.Inciso,
			Titulo:         d.Titulo,
			DescricaoLegal: d.DescricaoLegal,
			Produtos:       d.Produtos,
			Condicoes:      d.Condicoes,
		})
	}
	reducoesJSON, _ := json.Marshal(possibleReducoes)

	// Prompt focado na comparação de descrições e escolha entre candidatos
	prompt := fmt.Sprintf(`Como especialista tributário, analise se o produto abaixo se enquadra em alguma das previsões de redução fornecidas.
Cada documento candidato é independente e suas condições não se relacionam, escolha a mais específica ou adequada de acordo com a descrição do produto e as previsões legais.

Produto: %s
NCM: %s
Descrição NCM: %s
Reduções Candidatas: %s

Responda APENAS JSON:
{
  "carga_tributaria": 0,
  "justificativa": "Sua justificativa aqui em no máximo 20 palavras"
}

O campo "carga_tributaria" deve ser um número e não uma string.`, req.Descricao, req.NCM, ncmDescricao, string(reducoesJSON))

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
		return nil, fmt.Errorf("Gemini retornou um texto vazio para a análise de redução (finishReason: %s)", candidate.FinishReason)
	}

	var llmResult domains.ReducaoValidacaoResponse
	if err := json.Unmarshal([]byte(rawText), &llmResult); err != nil {
		return nil, fmt.Errorf("erro ao processar decisão da IA: %w (finishReason: %s, raw: %s)", err, candidate.FinishReason, rawText)
	}

	return &llmResult, nil
}

func (s *ReducaoService) CalcularReducao(ctx context.Context, req domains.ValidateRequest, reducaoResult *domains.ReducaoValidacaoResponse, aliquotaResult domains.AliquotaValidacaoResponse) (*domains.ReducaoValidacaoResponse, error) {
	if reducaoResult == nil || reducaoResult.CargaTributaria == 0 {
		return &domains.ReducaoValidacaoResponse{
			ReducaoBase:   0,
			Justificativa: "Produto não possui reduções previstas na legislação.",
		}, nil
	}

	if aliquotaResult.Aliquota == 0 {
		return &domains.ReducaoValidacaoResponse{
			ReducaoBase:     0,
			CargaTributaria: reducaoResult.CargaTributaria,
			Justificativa:   "Não foi possível calcular o percentual de redução pois a alíquota informada é zero.",
		}, nil
	}

	// Percentual de Redução da base de cálculo do ICMS
	// O cálculo é: (1 - (carga / aliquota)) * 100
	reducaoBase := (1 - (reducaoResult.CargaTributaria / aliquotaResult.Aliquota)) * 100
	reducaoBase = math.Round(reducaoBase*100) / 100

	return &domains.ReducaoValidacaoResponse{
		ReducaoBase:     reducaoBase,
		CargaTributaria: reducaoResult.CargaTributaria,
		Justificativa:   reducaoResult.Justificativa,
	}, nil
}
