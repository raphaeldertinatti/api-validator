package domains

// IcmsStDocument representa o documento de ICMS ST no MongoDB
type IcmsStDocument struct {
	ID              interface{}         `json:"id" bson:"_id,omitempty"`
	UF              string              `json:"uf" bson:"uf"`
	NCM             string              `json:"ncm" bson:"ncm"`
	TipoNCM         string              `json:"tipo_ncm" bson:"tipo_ncm"`
	CEST            string              `json:"cest" bson:"cest"`
	Descricao       string              `json:"descricao" bson:"descricao"`
	MVAOriginal     *float64            `json:"mva_original" bson:"mva_original"`
	TipoBaseCalculo string              `json:"tipo_base_calculo" bson:"tipo_base_calculo"`
	Fundamentacao   IcmsStFundamentacao `json:"fundamentacao" bson:"fundamentacao"`
	Vigencia        IcmsStVigencia      `json:"vigencia" bson:"vigencia"`
	Ativo           bool                `json:"ativo" bson:"ativo"`
}

// IcmsStFundamentacao representa as fundamentações legais do ICMS ST
type IcmsStFundamentacao struct {
	Portaria    string `json:"portaria" bson:"portaria"`
	ArtigoRICMS string `json:"artigo_ricms" bson:"artigo_ricms"`
	Convenio    string `json:"convenio" bson:"convenio"`
}

// IcmsStVigencia representa o período de vigência das regras do ICMS ST
type IcmsStVigencia struct {
	Inicio string  `json:"inicio" bson:"inicio"`
	Fim    *string `json:"fim" bson:"fim"`
}

type IcmsStValidacaoResponse struct {
	MVA             float64 `json:"mva,omitempty"`
	Portaria        string  `json:"portaria,omitempty"`
	TipoBaseCalculo string  `json:"tipo_base_calculo,omitempty"`
}
