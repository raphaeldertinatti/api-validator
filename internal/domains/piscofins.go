package domains

type PISCOFINSDocument struct {
	ID              interface{}         `bson:"_id" json:"id"`
	NCM             string              `bson:"ncm" json:"ncm"`
	CodigoFormatado string              `bson:"codigo_formatado" json:"codigo_formatado"`
	TipoNCM         string              `bson:"tipo_ncm" json:"tipo_ncm"`
	Regime          string              `bson:"regime" json:"regime"`
	CST             PISCOFINSCST        `bson:"cst" json:"cst"`
	Descricao       string              `bson:"descricao_produto" json:"descricao_produto"`
	Grupo           string              `bson:"grupo" json:"grupo"`
	Fundamentacao   PISCOFINSFundamento `bson:"fundamentacao" json:"fundamentacao"`
	Ativo           bool                `bson:"ativo" json:"ativo"`
}

type PISCOFINSCST struct {
	Entrada string `bson:"entrada" json:"entrada"`
	Saida   string `bson:"saida" json:"saida"`
}

type PISCOFINSFundamento struct {
	TabelaSped string `bson:"tabela_sped" json:"tabela_sped,omitempty"`
}

type PISCOFINSValidacaoResponse struct {
	Regime           string `json:"regime,omitempty"`
	Descricao        string `json:"descricao_produto,omitempty"`
	Grupo            string `json:"grupo,omitempty"`
	CSTPisEntrada    string `json:"cst_pis_entrada,omitempty"`
	CSTPisSaida      string `json:"cst_pis_saida,omitempty"`
	CSTCofinsEntrada string `json:"cst_cofins_entrada,omitempty"`
	CSTCofinsSaida   string `json:"cst_cofins_saida,omitempty"`
}
