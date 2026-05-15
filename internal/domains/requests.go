package domains

type ValidateRequest struct {
	NCM               string `json:"ncm" binding:"required"`
	Descricao         string `json:"descricao" binding:"required"`
	FornecedorSimples bool   `json:"fornecedor_simples"`
}
