package domains

type CESTDocument struct {
	ID          interface{}  `bson:"_id" json:"id"`
	CEST        string       `bson:"cest" json:"cest"`
	CESTLimpo   string       `bson:"cest_limpo" json:"cest_limpo"`
	Descricao   string       `bson:"descricao" json:"descricao"`
	NCMs        []CESTNCM    `bson:"ncms" json:"ncms"`
	NCMsCodigos []string     `bson:"ncms_codigos" json:"ncms_codigos"`
	Convenio    string       `bson:"convenio" json:"convenio"`
	Ativo       bool         `bson:"ativo" json:"ativo"`
	Segmento    CESTSegmento `bson:"segmento" json:"segmento"`
}

type CESTNCM struct {
	Codigo      string `bson:"codigo" json:"codigo"`
	CodigoLimpo string `bson:"codigo_limpo" json:"codigo_limpo"`
	Tipo        string `bson:"tipo" json:"tipo"`
}

type CESTSegmento struct {
	Anexo     int    `bson:"anexo" json:"anexo"`
	Descricao string `bson:"descricao" json:"descricao"`
}

type CESTValidacaoResponse struct {
	Enquadrado     *CESTResult  `json:"enquadrado,omitempty"`
	Possibilidades []CESTResult `json:"possibilidades,omitempty"`
	Justificativa  string       `json:"justificativa"`
	Status         string       `json:"status"` // DEFINIDO, AMBIGUO, NAO_ENQUADRADO
}

type CESTResult struct {
	CEST      string `json:"cest"`
	Descricao string `json:"descricao"`
	Segmento  string `json:"segmento"`
}
