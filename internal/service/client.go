package service

import (
	"fmt"
	"golang.org/x/text/encoding/htmlindex"
	"io"
	"net/http"
	"strings"
)

type HttpClient interface {
	Get(url string) (*http.Response, error)
}

type ClientService struct {
	hc HttpClient
}

func NewClientService(hc HttpClient) *ClientService {
	return &ClientService{
		hc: hc,
	}
}

func (c *ClientService) Do(url string) (string, error) {
	resp, err := c.hc.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка при запросе: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return "", err
	}

	charset := getResponseCharset(resp)

	if charset != "utf-8" {
		bodyBytes, err = decodeCharset(bodyBytes, charset)
		if err != nil {
			return "", err
		}
	}

	return string(bodyBytes), nil
}

func getResponseCharset(resp *http.Response) string {
	contentType := resp.Header.Get("Content-Type")
	charsetStart := strings.Index(contentType, "charset=")

	var charset string
	if charsetStart != -1 {
		charset = strings.TrimPrefix(contentType[charsetStart+len("charset="):], " ")
	} else {
		charset = "utf-8"
	}

	return charset
}

func decodeCharset(data []byte, charset string) ([]byte, error) {
	dec, err := htmlindex.Get(charset)
	if err != nil {
		return nil, err
	}

	utf8Data, err := dec.NewDecoder().Bytes(data)
	if err != nil {
		return nil, err
	}

	return utf8Data, nil
}
