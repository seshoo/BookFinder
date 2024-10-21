package service_test

import (
	"errors"
	"testing"

	"github.com/seshoo/bookFinder/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Do(url string) (string, error) {
	args := m.Called(url)
	return args.String(0), args.Error(1)
}

type MockExtractor struct {
	mock.Mock
}

func (m *MockExtractor) Get(content string) (string, string, error) {
	args := m.Called(content)
	return args.String(0), args.String(1), args.Error(2)
}

func TestParserService_Get(t *testing.T) {
	mockClient := new(MockClient)
	mockExtractor := new(MockExtractor)

	urlTemplate := "http://example.com/%s"
	parserService := service.NewParserService(urlTemplate, mockClient, mockExtractor)

	code := "test-code"
	expectedContent := "test-content"
	expectedTitle := "Test Title"
	expectedText := "Test Text"
	expectedPageData := service.PageData{Code: code, Title: expectedTitle, Text: expectedText}

	mockClient.On("Do", "http://example.com/test-code").Return(expectedContent, nil)
	mockExtractor.On("Get", expectedContent).Return(expectedTitle, expectedText, nil)

	pageData, err := parserService.Get(code)

	assert.NoError(t, err)
	assert.Equal(t, expectedPageData, pageData)

	mockClient.AssertExpectations(t)
	mockExtractor.AssertExpectations(t)
}

func TestParserService_Get_ClientError(t *testing.T) {
	mockClient := new(MockClient)
	mockExtractor := new(MockExtractor)

	urlTemplate := "http://example.com/%s"
	parserService := service.NewParserService(urlTemplate, mockClient, mockExtractor)

	code := "test-code"
	expectedError := errors.New("client error")

	mockClient.On("Do", "http://example.com/test-code").Return("", expectedError)

	pageData, err := parserService.Get(code)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, service.PageData{}, pageData)

	mockClient.AssertExpectations(t)
	mockExtractor.AssertExpectations(t)
}

func TestParserService_Get_ExtractorError(t *testing.T) {
	mockClient := new(MockClient)
	mockExtractor := new(MockExtractor)

	urlTemplate := "http://example.com/%s"
	parserService := service.NewParserService(urlTemplate, mockClient, mockExtractor)

	code := "test-code"
	expectedContent := "test-content"
	expectedError := errors.New("extractor error")

	mockClient.On("Do", "http://example.com/test-code").Return(expectedContent, nil)
	mockExtractor.On("Get", expectedContent).Return("", "", expectedError)

	pageData, err := parserService.Get(code)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, service.PageData{}, pageData)

	mockClient.AssertExpectations(t)
	mockExtractor.AssertExpectations(t)
}
