package domains

type ValidateResponse struct {
	NCM       *NCMValidacaoResponse       `json:"ncm,omitempty"`
	IPI       *IPIValidacaoResponse       `json:"ipi,omitempty"`
	CEST      *CESTValidacaoResponse      `json:"cest,omitempty"`
	PISCOFINS *PISCOFINSValidacaoResponse `json:"piscofins,omitempty"`
}
