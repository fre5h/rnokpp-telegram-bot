package mocks

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type MockOkTelegramClient struct{}

func NewMockOkClient() *MockOkTelegramClient {
	return &MockOkTelegramClient{}
}

func (c MockOkTelegramClient) SendTextMessageToChat(_ int, _ string) (string, error) {
	return "OK", nil
}

type MockFailedTelegramClient struct{}

func NewMockFailedClient() *MockFailedTelegramClient {
	return &MockFailedTelegramClient{}
}

func (c MockFailedTelegramClient) SendTextMessageToChat(_ int, _ string) (string, error) {
	return "FAIL", errors.New("failed request")
}

type MockHttpClient struct {
	statusCode int
	body       string
	err        error
}

func NewMockHttpClient(statusCode int, body string, err error) *MockHttpClient {
	return &MockHttpClient{statusCode, body, err}
}

func (c MockHttpClient) PostForm(_ string, _ url.Values) (*http.Response, error) {
	return &http.Response{
		StatusCode: c.statusCode,
		Body:       io.NopCloser(strings.NewReader(c.body)),
	}, c.err
}

type MockHttpClientFailedCloser struct {
	statusCode int
	body       string
	err        error
}

func NewMockHttpClientFailedCloser(statusCode int, body string, err error) *MockHttpClientFailedCloser {
	return &MockHttpClientFailedCloser{statusCode, body, err}
}

func (c MockHttpClientFailedCloser) PostForm(_ string, _ url.Values) (*http.Response, error) {
	return &http.Response{
		StatusCode: c.statusCode,
		Body: &MockReadCloser{
			reader: strings.NewReader(c.body),
		},
	}, nil
}

type MockReadCloser struct {
	reader io.Reader
}

func (r MockReadCloser) Read(_ []byte) (n int, err error) {
	return 0, errors.New("error on read")
}

func (r MockReadCloser) Close() error {
	return fmt.Errorf("error on close")
}
