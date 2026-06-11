package domains

type ValidateResponse struct {
	NCM         *NCMValidacaoResponse         `json:"ncm,omitempty"`
	IPI         *IPIValidacaoResponse         `json:"ipi,omitempty"`
	CEST        *CESTValidacaoResponse        `json:"cest,omitempty"`
	PISCOFINS   *PISCOFINSValidacaoResponse   `json:"piscofins,omitempty"`
	Isencao     *IsencaoValidacaoResponse     `json:"isencao,omitempty"`
	Diferimento *DiferimentoValidacaoResponse `json:"diferimento,omitempty"`
	Aliquota    *AliquotaValidacaoResponse    `json:"aliquota,omitempty"`
	Reducao     *ReducaoValidacaoResponse     `json:"reducao,omitempty"`
}
