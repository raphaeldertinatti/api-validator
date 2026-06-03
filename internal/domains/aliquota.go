package domains

type AliquotaDocument struct {
	ID             interface{}           `bson:"_id" json:"id"`
	UF             string                `bson:"uf" json:"uf"`
	Artigo         string                `bson:"artigo" json:"artigo"`
	Inciso         string                `bson:"inciso" json:"inciso"`
	Titulo         string                `bson:"titulo" json:"titulo"`
	Tipo           string                `bson:"tipo" json:"tipo"`
	Subtipo        string                `bson:"subtipo" json:"subtipo"`
	DescricaoLegal string                `bson:"descricao_legal" json:"descricao_legal"`
	Aliquota       float64               `bson:"aliquota" json:"aliquota"`
	Produtos       []string              `bson:"produtos" json:"produtos"`
	NCMS           []AliquotaNCM         `bson:"ncms" json:"ncms"`
	NCMSCodigos    []string              `bson:"ncms_codigos" json:"ncms_codigos"`
	CapitulosNCM   []string              `bson:"capitulos_ncm" json:"capitulos_ncm"`
	Condicoes      AliquotaCondicoes     `bson:"condicoes" json:"condicoes"`
	Observacoes    []string              `bson:"observacoes" json:"observacoes"`
	Fundamentacao  AliquotaFundamentacao `bson:"fundamentacao" json:"fundamentacao"`
	CSTICMS        string                `bson:"cst_icms" json:"cst_icms"`
	Vigencia       AliquotaVigencia      `bson:"vigencia" json:"vigencia"`
	Ativo          bool                  `bson:"ativo" json:"ativo"`
}

type AliquotaNCM struct {
	Codigo  string `bson:"codigo" json:"codigo"`
	Tipo    string `bson:"tipo" json:"tipo"`
	Produto string `bson:"produto" json:"produto"`
}

type AliquotaCondicoes struct {
	Incluir []string `bson:"incluir" json:"incluir"`
	Excluir []string `bson:"excluir" json:"excluir"`
}

type AliquotaFundamentacao struct {
	Artigo       string `bson:"artigo" json:"artigo"`
	Inciso       string `bson:"inciso" json:"inciso"`
	Lei          string `bson:"lei" json:"lei"`
	DecretoAtual string `bson:"decreto_atual" json:"decreto_atual"`
	RICMS        string `bson:"ricms" json:"ricms"`
}

type AliquotaVigencia struct {
	Inicio string  `bson:"inicio" json:"inicio"`
	Fim    *string `bson:"fim" json:"fim"`
}

type AliquotaValidacaoResponse struct {
	Aliquota      float64 `json:"aliquota"`
	Artigo        string  `json:"artigo,omitempty"`
	Justificativa string  `json:"justificativa,omitempty"`
	Observacao    string  `json:"observacao,omitempty"`
}
