package domains

type ValidateRequest struct {
	NCM               string `json:"ncm" binding:"required"`
	Descricao         string `json:"descricao" binding:"required"`
	FornecedorSimples bool   `json:"fornecedor_simples"`
}

type NCMHierarchyItem struct {
	Codigo    string `bson:"codigo"`
	Descricao string `bson:"descricao"`
}

type ValidateResponse struct {
	NCMDescricao    string `json:"ncm_descricao"`
	ValidacaoLLM    string `json:"validacao_llm"`
	StatusValidacao bool   `json:"status_validacao"`
}
