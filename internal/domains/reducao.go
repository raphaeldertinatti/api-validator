package domains

type ReducaoDocument struct {
	ID              interface{}          `bson:"_id,omitempty" json:"id"`
	UF              string               `bson:"uf" json:"uf"`
	Artigo          string               `bson:"artigo" json:"artigo"`
	Inciso          string               `bson:"inciso" json:"inciso"`
	Titulo          string               `bson:"titulo" json:"titulo"`
	Tipo            string               `bson:"tipo" json:"tipo"`
	Subtipo         string               `bson:"subtipo" json:"subtipo"`
	DescricaoLegal  string               `bson:"descricao_legal" json:"descricao_legal"`
	CargaTributaria float64              `bson:"carga_tributaria" json:"carga_tributaria"`
	Produtos        []string             `bson:"produtos" json:"produtos"`
	NCMS            []ReducaoNCM         `bson:"ncms" json:"ncms"`
	NCMSCodigos     []string             `bson:"ncms_codigos" json:"ncms_codigos"`
	CapitulosNCM    []string             `bson:"capitulos_ncm" json:"capitulos_ncm"`
	Condicoes       ReducaoCondicoes     `bson:"condicoes" json:"condicoes"`
	Observacoes     []string             `bson:"observacoes" json:"observacoes"`
	Convenios       []string             `bson:"convenios" json:"convenios"`
	Fundamentacao   ReducaoFundamentacao `bson:"fundamentacao" json:"fundamentacao"`
	CSTICMS         string               `bson:"cst_icms" json:"cst_icms"`
	Vigencia        ReducaoVigencia      `bson:"vigencia" json:"vigencia"`
	Ativo           bool                 `bson:"ativo" json:"ativo"`
}

type ReducaoNCM struct {
	Codigo  string `bson:"codigo" json:"codigo"`
	Tipo    string `bson:"tipo" json:"tipo"`
	Produto string `bson:"produto" json:"produto"`
}

type ReducaoCondicoes struct {
	Incluir []string `bson:"incluir" json:"incluir"`
	Excluir []string `bson:"excluir" json:"excluir"`
}

type ReducaoFundamentacao struct {
	Artigo       string `bson:"artigo" json:"artigo"`
	Inciso       string `bson:"inciso" json:"inciso"`
	Anexo        string `bson:"anexo" json:"anexo"`
	DecretoAtual string `bson:"decreto_atual" json:"decreto_atual"`
	RICMS        string `bson:"ricms" json:"ricms"`
}

type ReducaoVigencia struct {
	Inicio string `bson:"inicio" json:"inicio"`
	Fim    string `bson:"fim" json:"fim"`
}

type ReducaoValidacaoResponse struct {
	CargaTributaria float64 `json:"carga_tributaria,omitempty"`
	ReducaoBase     float64 `json:"reducao_base"`
	Justificativa   string  `json:"justificativa,omitempty"`
}
