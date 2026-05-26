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

type IPIRepository interface {
	FindByCode(ctx context.Context, ncmCode string) (*domains.IPIDocument, error)
}

type IPIService struct {
	repo   IPIRepository
	gemini *GeminiService
}

func NewIPIService(repo IPIRepository, gemini *GeminiService) *IPIService {
	return &IPIService{
		repo:   repo,
		gemini: gemini,
	}
}

func (s *IPIService) ValidateIPI(ctx context.Context, req domains.ValidateRequest) (*domains.IPIValidacaoResponse, error) {
	// 1. Buscar dados do IPI no MongoDB
	doc, err := s.repo.FindByCode(ctx, req.NCM)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar IPI: %w", err)
	}
	if doc == nil {
		return &domains.IPIValidacaoResponse{Situacao: "NCM não encontrada"}, nil
	}

	// 2. Se tiver EX Tarifário, chama a lógica especializada
	if doc.TemEx && len(doc.ExTarifarios) > 0 {
		return s.EnquadramentoExIPI(ctx, req, doc)
	}

	return &domains.IPIValidacaoResponse{
		Aliquota: doc.Aliquota,
		Situacao: doc.Situacao,
	}, nil
}

func (s *IPIService) EnquadramentoExIPI(ctx context.Context, req domains.ValidateRequest, doc *domains.IPIDocument) (*domains.IPIValidacaoResponse, error) {
	exsJSON, _ := json.Marshal(doc.ExTarifarios)

	prompt := fmt.Sprintf(`Você é um especialista tributário. Classifique o Produto nos EX da NCM %s.
Prod: %s
EXs: %s

Responda APENAS JSON:
{"status":"DEFINIDO"|"AMBIGUO"|"PADRAO","ex_enquadrado":{"ex":"","aliquota":0,"descricao":"","justificativa":""},"possibilidades":[],"justificativa":"Sua justificativa aqui em no máximo 20 palavras"}

DEFINIDO: 1 match claro. AMBIGUO: dúvida (liste em possibilidades). PADRAO: 0 match.`, req.NCM, req.Descricao, string(exsJSON))

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

	rawText := strings.TrimSpace(gemResp.Candidates[0].Content.Parts[0].Text)
	rawText = strings.TrimPrefix(rawText, "```json")
	rawText = strings.TrimPrefix(rawText, "```")
	rawText = strings.TrimSuffix(rawText, "```")
	rawText = strings.TrimSpace(rawText)

	if rawText == "" {
		return nil, fmt.Errorf("Gemini retornou um texto vazio para enquadramento de IPI")
	}

	var llmResult struct {
		Status         string                      `json:"status"`
		ExEnquadrado   *domains.ExTarifarioResult  `json:"ex_enquadrado"`
		Possibilidades []domains.ExTarifarioResult `json:"possibilidades"`
		Justificativa  string                      `json:"justificativa"`
	}

	if err := json.Unmarshal([]byte(rawText), &llmResult); err != nil {
		return nil, fmt.Errorf("erro ao processar decisão da IA: %w (raw: %s)", err, rawText)
	}

	finalResp := &domains.IPIValidacaoResponse{
		Aliquota:      doc.Aliquota,
		Situacao:      doc.Situacao,
		Justificativa: llmResult.Justificativa,
	}

	if llmResult.Status == "DEFINIDO" && llmResult.ExEnquadrado != nil {
		finalResp.ExEnquadrado = llmResult.ExEnquadrado
		finalResp.Aliquota = llmResult.ExEnquadrado.Aliquota
	} else if llmResult.Status == "AMBIGUO" {
		finalResp.Possibilidades = llmResult.Possibilidades
		finalResp.Justificativa = "Não foi possível definir o EX por falta de informações: " + llmResult.Justificativa
	}

	return finalResp, nil
}
