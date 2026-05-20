package domains

type IPIDocument struct {
	ID              interface{} `bson:"_id" json:"id"`
	NCM             string      `bson:"ncm" json:"ncm"`
	CodigoFormatado string      `bson:"codigo_formatado" json:"codigo_formatado"`
	Descricao       string      `bson:"descricao" json:"descricao"`
	Aliquota        *float64    `bson:"aliquota" json:"aliquota"`
	Situacao        string         `bson:"situacao" json:"situacao"`
	ExTarifarios    []ExTarifario  `bson:"ex_tarifarios" json:"ex_tarifarios"`
	TemEx           bool           `bson:"tem_ex" json:"tem_ex"`
	Vigencia        struct {
		Inicio string  `bson:"inicio" json:"inicio"`
		Fim    *string `bson:"fim" json:"fim"`
	} `bson:"vigencia" json:"vigencia"`
	Fundamentacao struct {
		DecretoBase string `bson:"decreto_base" json:"decreto_base"`
	} `bson:"fundamentacao" json:"fundamentacao"`
	Versao string `bson:"versao" json:"versao"`
	Ativo  bool   `bson:"ativo" json:"ativo"`
}

type ExTarifario struct {
	Ex        string   `bson:"ex" json:"ex"`
	Descricao string   `bson:"descricao" json:"descricao"`
	Aliquota  *float64 `bson:"aliquota" json:"aliquota"`
	Situacao  string   `bson:"situacao" json:"situacao"`
}

type IPIValidacaoResponse struct {
	Aliquota      *float64              `json:"aliquota"`
	Situacao      string                `json:"situacao"`
	ExEnquadrado  *ExTarifarioResult    `json:"ex_enquadrado,omitempty"`
	Possibilidades []ExTarifarioResult   `json:"possibilidades,omitempty"`
	Justificativa string                `json:"justificativa,omitempty"`
}

type ExTarifarioResult struct {
	Ex            string   `json:"ex"`
	Aliquota      *float64 `json:"aliquota"`
	Descricao     string   `json:"descricao"`
	Justificativa string   `json:"justificativa,omitempty"`
}
