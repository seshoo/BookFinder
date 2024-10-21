package service

import "net/http"

type PageData struct {
	Code  string
	Title string
	Text  string
}

type Parser interface {
	Get(code string) (PageData, error)
}

type Services struct {
	Parser Parser
}

type Deps struct {
	DpUrlTmp string
}

func NewServices(deps Deps) *Services {
	clientService := NewClientService(http.DefaultClient)
	extractorService := NewExtractorService()

	return &Services{
		Parser: NewParserService(deps.DpUrlTmp, clientService, extractorService),
	}
}
