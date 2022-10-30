package handler

import (
	"net/http"
	"net/url"
)

type HttpClient interface {
	PostForm(string, url.Values) (*http.Response, error)
}
