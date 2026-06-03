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

type AliquotaRepository interface {
	FindByCode(ctx context.Context, ncmCode string) ([]domains.AliquotaDocument, error)
}

type AliquotaService struct {
	ncmRepo NCMRepository
	repo    AliquotaRepository
	gemini  *GeminiService
}

func NewAliquotaService(ncmRepo NCMRepository, repo AliquotaRepository, gemini *GeminiService) *AliquotaService {
	return &AliquotaService{
		ncmRepo: ncmRepo,
		repo:    repo,
		gemini:  gemini,
	}
}

func (s *AliquotaService) ValidateAliquota(ctx context.Context, req domains.ValidateRequest) (*domains.AliquotaValidacaoResponse, error) {
	// 1. Buscar possíveis aliquotas para a NCM no MongoDB
	docs, err := s.repo.FindByCode(ctx, req.NCM)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar aliquotas: %w", err)
	}

	ncmDescricao, err := s.ncmRepo.FindDescricaoCompleta(ctx, req.NCM)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar NCM: %w", err)
	}

	if len(docs) == 0 { //TODO: Quando tiver mais estados implementar lógica de fallback para outros estados conforme UF do request
		return &domains.AliquotaValidacaoResponse{
			Aliquota:      18,
			Justificativa: "Aliquota Padrão",
		}, nil
	}

	// 2. Analisar se a descrição do produto bate com a descrição de alguma aliquota
	return s.AnaliseAliquota(ctx, req, docs, ncmDescricao)
}

func (s *AliquotaService) AnaliseAliquota(ctx context.Context, req domains.ValidateRequest, docs []domains.AliquotaDocument, ncmDescricao string) (*domains.AliquotaValidacaoResponse, error) {
	type AliquotaBrief struct {
		ID             interface{}               `json:"id"`
		Aliquota       float64                   `json:"aliquota"`
		Artigo         string                    `json:"artigo"`
		Inciso         string                    `json:"inciso"`
		Titulo         string                    `json:"titulo"`
		DescricaoLegal string                    `json:"descricao_legal"`
		NCMS           []string                  `json:"ncm_codigos,omitempty"`
		CapitulosNCM   []string                  `json:"capitulos_ncm"`
		Produtos       []string                  `json:"produtos"`
		Condicoes      domains.AliquotaCondicoes `json:"condicoes"`
	}
	var possibleAliquota []AliquotaBrief
	for _, d := range docs {
		possibleAliquota = append(possibleAliquota, AliquotaBrief{
			ID:             d.ID,
			Aliquota:       d.Aliquota,
			Artigo:         d.Artigo,
			Inciso:         d.Inciso,
			Titulo:         d.Titulo,
			DescricaoLegal: d.DescricaoLegal,
			NCMS:           d.NCMSCodigos,
			CapitulosNCM:   d.CapitulosNCM,
			Produtos:       d.Produtos,
			Condicoes:      d.Condicoes,
		})
	}
	aliquotasJSON, _ := json.Marshal(possibleAliquota)

	// Prompt focado na comparação de descrições e escolha entre candidatos
	prompt := fmt.Sprintf(`Como especialista tributário, analise se o produto abaixo se enquadra em alguma das previsões de aliquota ICMS fornecidas.
Cada documento candidato é independente e suas condições não se relacionam, escolha a mais específica ou adequada de acordo com a descrição do produto e as previsões legais.

Produto: %s
NCM: %s
Descrição NCM: %s
Aliquotas Candidatas: %s

Responda APENAS JSON:
{
  "aliquota": 0,
  "artigo": "...",
  "justificativa": "Sua justificativa aqui em no máximo 20 palavras"
}

O campo "aliquota" deve ser um float64 e não uma string.
Caso nenhum candidato se aplique alíquota será 18, a justificativa será 'Aliquota Padrão' e omita o campo artigo`, req.Descricao, req.NCM, ncmDescricao, string(aliquotasJSON))

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
		return nil, fmt.Errorf("Gemini retornou um texto vazio para a análise de alíquota (finishReason: %s)", candidate.FinishReason)
	}

	var llmResult domains.AliquotaValidacaoResponse
	if err := json.Unmarshal([]byte(rawText), &llmResult); err != nil {
		return nil, fmt.Errorf("erro ao processar decisão da IA: %w (finishReason: %s, raw: %s)", err, candidate.FinishReason, rawText)
	}

	return &llmResult, nil
}
