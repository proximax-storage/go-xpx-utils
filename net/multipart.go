package net

import (
	"context"
	"errors"
	"io"
	"net/http"
)

type MultiPartHttpClient struct {
	addr string
	cl   *http.Client
}

func NewMultiPartHttpClientClient(addr string) (*MultiPartHttpClient, error) {
	if len(addr) == 0 {
		return nil, errors.New("address should not be blank")
	}

	return &MultiPartHttpClient{addr: addr, cl: &http.Client{}}, nil
}

func (ref *MultiPartHttpClient) PostFile(ctx context.Context, path string, reqBody io.Reader, headerRaws ...HeaderRaw) (io.Reader, HttpStatusCode, error) {
	req, err := http.NewRequest(http.MethodPost, ref.addr+path, reqBody)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "multipart/form-data")

	for _, headerRaw := range headerRaws {
		req.Header.Set(headerRaw.key, headerRaw.value)
	}

	req.WithContext(ctx)

	resp, err := ref.cl.Do(req)
	if err != nil {
		return nil, 0, err
	}

	return resp.Body, HttpStatusCode(resp.StatusCode), nil
}

func (ref *MultiPartHttpClient) GetFile(ctx context.Context, path string, headerRaws ...HeaderRaw) (io.Reader, HttpStatusCode, error) {
	req, err := http.NewRequest(http.MethodPost, ref.addr+path, nil)
	if err != nil {
		return nil, 0, err
	}

	for _, headerRaw := range headerRaws {
		req.Header.Set(headerRaw.key, headerRaw.value)
	}

	req.WithContext(ctx)

	resp, err := ref.cl.Do(req)
	if err != nil {
		return nil, 0, err
	}

	return resp.Body, HttpStatusCode(resp.StatusCode), nil
}
