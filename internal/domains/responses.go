package domains

type ValidateResponse struct {
	NCM *NCMValidacaoResponse `json:"ncm,omitempty"`
	IPI *IPIValidacaoResponse `json:"ipi,omitempty"`
}
