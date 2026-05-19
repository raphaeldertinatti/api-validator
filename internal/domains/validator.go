package domains

type ValidationResult struct {
	NCM NCMValidacaoResponse `json:"ncm"`
	// IPI       IPIResult       `json:"ipi"`
	// CEST      CESTResult      `json:"cest"`
	// PISCofins PISCOFINSResult `json:"pis_cofins"`
	// ICMS      ICMSResult      `json:"icms"`
}
