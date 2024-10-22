package service

import (
	"github.com/seshoo/bookFinder/internal/repository"
	"net/http"
)

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
	Repositories repository.Repositories
	DpUrlTmp     string
}

func NewServices(deps Deps) *Services {
	clientService := NewClientService(http.DefaultClient)
	extractorService := NewExtractorService()

	return &Services{
		Parser: NewParserService(deps.DpUrlTmp, clientService, extractorService),
	}
}
