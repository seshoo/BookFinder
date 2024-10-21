package service

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type ExtractorService struct {
}

func NewExtractorService() *ExtractorService {
	return &ExtractorService{}
}

func (e *ExtractorService) Get(content string) (string, string, error) {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(content))

	return doc.Find("Title").Text(), doc.Find("div.sp-body").Text(), nil
}
