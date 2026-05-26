package domains

import "fmt"

// BuildNCMCandidates gera prefixos de NCM baseados nos comprimentos fornecidos.
// Ex: ("85442000", []int{8, 4, 2}) -> ["85442000", "8544", "85"]
func BuildNCMCandidates(ncm string, lengths []int) []string {
	seen := make(map[string]struct{})
	var candidates []string

	for _, l := range lengths {
		if len(ncm) >= l {
			prefix := ncm[:l]
			if _, exists := seen[prefix]; !exists {
				seen[prefix] = struct{}{}
				candidates = append(candidates, prefix)
			}
		}
	}
	return candidates
}

// DocKey retorna a chave única de identificação para IsencaoDocument
func (d IsencaoDocument) DocKey() string {
	return fmt.Sprintf("%s_%s_%s_%s",
		d.Artigo,
		d.Paragrafo,
		d.Inciso,
		d.Subtipo,
	)
}

// DocKey retorna a chave única de identificação para DiferimentoDocument
func (d DiferimentoDocument) DocKey() string {
	inciso := ""
	if d.Inciso != nil {
		inciso = *d.Inciso
	}
	return fmt.Sprintf("%s_%s_%s",
		d.Artigo,
		inciso,
		d.Titulo,
	)
}
