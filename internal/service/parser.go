package service

import "fmt"

type Client interface {
	Do(url string) (string, error)
}

type Extractor interface {
	Get(content string) (string, string, error)
}

type ParserService struct {
	urlTemplate string
	client      Client
	extractor   Extractor
}

func NewParserService(urlTemplate string, client Client, extractor Extractor) *ParserService {
	return &ParserService{
		urlTemplate: urlTemplate,
		client:      client,
		extractor:   extractor,
	}
}

func (p *ParserService) Get(code string) (PageData, error) {
	var pd PageData

	url := fmt.Sprintf(p.urlTemplate, code)
	content, err := p.client.Do(url)
	if err != nil {
		return PageData{}, err
	}

	title, text, err := p.extractor.Get(content)
	if err != nil {
		return PageData{}, err
	}

	pd.Code = code
	pd.Title = title
	pd.Text = text

	return pd, nil
}
