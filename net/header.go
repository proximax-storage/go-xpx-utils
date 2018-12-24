package net

import "net/http"

type RequestOption func(req *http.Request)

func NewHeaderRow(key, value string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set(key, value)
	}
}
