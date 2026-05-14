package domains

type NCMDocument struct {
	ID              interface{} `bson:"_id"`
	NCM             string      `bson:"ncm"`
	CodigoFormatado string      `bson:"codigo_formatado"`
	Descricao       struct {
		Curta       string `bson:"curta"`
		Completa    string `bson:"completa"`
		Normalizada string `bson:"normalizada"`
	} `bson:"descricao"`
	Hierarquia struct {
		Capitulo   NCMHierarchyItem `bson:"capitulo"`
		Posicao    NCMHierarchyItem `bson:"posicao"`
		Subposicao NCMHierarchyItem `bson:"subposicao"`
		Subitem    NCMHierarchyItem `bson:"subitem"`
	} `bson:"hierarquia"`
	NCMPai        string   `bson:"ncm_pai"`
	PalavrasChave []string `bson:"palavras_chave"`
	Sinonimos     []string `bson:"sinonimos"`
	Ativo         bool     `bson:"ativo"`
	Vigencia      struct {
		Inicio string `bson:"inicio"`
		Fim    string `bson:"fim"`
	} `bson:"vigencia"`
	Ato struct {
		Tipo   string `bson:"tipo"`
		Numero string `bson:"numero"`
		Ano    string `bson:"ano"`
	} `bson:"ato"`
	Versao string `bson:"versao"`
}
