package domains

type DiferimentoDocument struct {
	ID                      interface{}              `json:"id" bson:"_id,omitempty"`
	UF                      string                   `json:"uf" bson:"uf"`
	Artigo                  string                   `json:"artigo" bson:"artigo"`
	Inciso                  *string                  `json:"inciso" bson:"inciso"`
	Titulo                  string                   `json:"titulo" bson:"titulo"`
	Tipo                    string                   `json:"tipo" bson:"tipo"`
	DescricaoLegal          string                   `json:"descricao_legal" bson:"descricao_legal"`
	Produtos                []string                 `json:"produtos" bson:"produtos"`
	NCMs                    []string                 `json:"ncms" bson:"ncms"`
	NCMsCodigos             []string                 `json:"ncms_codigos" bson:"ncms_codigos"`
	CapitulosNCM            []string                 `json:"capitulos_ncm" bson:"capitulos_ncm"`
	Condicoes               DiferimentoCondicoes     `json:"condicoes" bson:"condicoes"`
	EncerramentoDiferimento string                   `json:"encerramento_diferimento" bson:"encerramento_diferimento"`
	Observacoes             []string                 `json:"observacoes" bson:"observacoes"`
	Fundamentacao           DiferimentoFundamentacao `json:"fundamentacao" bson:"fundamentacao"`
	CSTICMS                 string                   `json:"cst_icms" bson:"cst_icms"`
	Vigencia                DiferimentoVigencia      `json:"vigencia" bson:"vigencia"`
	Ativo                   bool                     `json:"ativo" bson:"ativo"`
}

type DiferimentoCondicoes struct {
	Incluir []string `json:"incluir" bson:"incluir"`
	Excluir []string `json:"excluir" bson:"excluir"`
}

type DiferimentoFundamentacao struct {
	Artigo       string `json:"artigo" bson:"artigo"`
	Lei          string `json:"lei" bson:"lei"`
	DecretoAtual string `json:"decreto_atual" bson:"decreto_atual"`
	RICMS        string `json:"ricms" bson:"ricms"`
}

type DiferimentoVigencia struct {
	Inicio string  `json:"inicio" bson:"inicio"`
	Fim    *string `json:"fim" bson:"fim"`
}

type DiferimentoValidacaoResponse struct {
	Diferido      bool   `json:"diferido"`
	Artigo        string `json:"artigo,omitempty"`
	Justificativa string `json:"justificativa,omitempty"`
	Observacao    string `json:"observacao,omitempty"`
}
