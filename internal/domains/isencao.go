package domains

type IsencaoDocument struct {
	ID             interface{}          `json:"id" bson:"_id,omitempty"`
	UF             string               `json:"uf" bson:"uf"`
	Artigo         string               `json:"artigo" bson:"artigo"`
	Paragrafo      string               `json:"paragrafo" bson:"paragrafo"`
	Inciso         string               `json:"inciso" bson:"inciso"`
	Titulo         string               `json:"titulo" bson:"titulo"`
	Tipo           string               `json:"tipo" bson:"tipo"`
	Subtipo        string               `json:"subtipo" bson:"subtipo"`
	DescricaoLegal string               `json:"descricao_legal" bson:"descricao_legal"`
	Produtos       []string             `json:"produtos" bson:"produtos"`
	NCMS           []IsencaoNCM         `json:"ncms" bson:"ncms"`
	NCMSCodigos    []string             `json:"ncms_codigos" bson:"ncms_codigos"`
	CapitulosNCM   []string             `json:"capitulos_ncm" bson:"capitulos_ncm"`
	Condicoes      IsencaoCondicoes     `json:"condicoes" bson:"condicoes"`
	Observacoes    []string             `json:"observacoes" bson:"observacoes"`
	Convenios      []string             `json:"convenios" bson:"convenios"`
	CSTICMS        string               `json:"cst_icms" bson:"cst_icms"`
	Vigencia       IsencaoVigencia      `json:"vigencia" bson:"vigencia"`
	Fundamentacao  IsencaoFundamentacao `json:"fundamentacao" bson:"fundamentacao"`
	Ativo          bool                 `json:"ativo" bson:"ativo"`
}

type IsencaoNCM struct {
	Codigo  string `json:"codigo" bson:"codigo"`
	Tipo    string `json:"tipo" bson:"tipo"`
	Produto string `json:"produto" bson:"produto"`
}

type IsencaoCondicoes struct {
	Incluir []string `json:"incluir" bson:"incluir"`
	Excluir []string `json:"excluir" bson:"excluir"`
}

type IsencaoVigencia struct {
	Inicio string `json:"inicio" bson:"inicio"`
	Fim    string `json:"fim" bson:"fim"`
}

type IsencaoFundamentacao struct {
	Artigo       string `json:"artigo" bson:"artigo"`
	Paragrafo    string `json:"paragrafo" bson:"paragrafo"`
	Inciso       string `json:"inciso" bson:"inciso"`
	Anexo        string `json:"anexo" bson:"anexo"`
	DecretoAtual string `json:"decreto_atual" bson:"decreto_atual"`
	RICMS        string `json:"ricms" bson:"ricms"`
}

type IsencaoValidacaoResponse struct {
	Isento        bool   `json:"isento"`
	Artigo        string `json:"artigo,omitempty"`
	Paragrafo     string `json:"paragrafo,omitempty"`
	Justificativa string `json:"justificativa,omitempty"`
}
