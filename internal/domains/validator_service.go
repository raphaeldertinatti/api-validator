package domains

type ValidateRequest struct {
	NCM               string `json:"ncm" binding:"required"`
	Descricao         string `json:"descricao" binding:"required"`
	FornecedorSimples bool   `json:"fornecedor_simples"`
}

type NCMDocument struct {
	ID              interface{} `bson:"_id"`
	NCM             string      `bson:"ncm"`
	CodigoFormatado string      `bson:"codigo_formatado"`
	Descricao       struct {
		Curta       string `bson:"curta"`
		Completa    string `bson:"completa"`
		Normalizada string `bson:"normalizada"`
	} `bson:"descricao"`
	Hierarquia struct {
		Capitulo   NCMHierarchyItem `bson:"capitulo"`
		Posicao    NCMHierarchyItem `bson:"posicao"`
		Subposicao NCMHierarchyItem `bson:"subposicao"`
		Subitem    NCMHierarchyItem `bson:"subitem"`
	} `bson:"hierarquia"`
	NCMPai        string   `bson:"ncm_pai"`
	PalavrasChave []string `bson:"palavras_chave"`
	Sinonimos     []string `bson:"sinonimos"`
	Ativo         bool     `bson:"ativo"`
	Vigencia      struct {
		Inicio string `bson:"inicio"`
		Fim    string `bson:"fim"`
	} `bson:"vigencia"`
	Ato struct {
		Tipo   string `bson:"tipo"`
		Numero string `bson:"numero"`
		Ano    string `bson:"ano"`
	} `bson:"ato"`
	Versao string `bson:"versao"`
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
