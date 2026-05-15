package domains

type ValidateResponse struct {
	NCMDescricao    string `json:"ncm_descricao"`
	ValidacaoLLM    string `json:"validacao_llm"`
	StatusValidacao bool   `json:"status_validacao"`
}
