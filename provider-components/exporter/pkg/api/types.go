package api

type EnvConfig struct {
	Port     string `default:"8081"`
	OdataUrl string `required:"true" envconfig:"odata_url"`
}

type CSVExporter interface {
	GetCSVBytes() []byte
}

type FaqResponse struct {
	Faqs *[]Faqs `json:"value"`
	CSVExporter
}

type Faqs struct {
	ID         string        `json:"ID"`
	Title      string        `json:"title"`
	Descr      string        `json:"descr"`
	State      string        `json:"state"`
	Answer     string        `json:"answer"`
	Author     string        `json:"author"`
	CategoryID int           `json:"category_ID"`
	Categories *[]Categories `json:"category"`
	Count      int           `json:"count"`
}

type Authors struct {
	ID   string  `json:"ID"`
	Name string  `json:"name"`
	Faqs *[]Faqs `json:"faqs"`
}

type Categories struct {
	ID   int    `json:"ID"`
	Name string `json:"name"`
}
