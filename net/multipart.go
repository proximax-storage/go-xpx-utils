package net

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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

func (ref *MultiPartHttpClient) PostFile(ctx context.Context, path string, fileParamName, filePath string, headerRaws ...HeaderRaw) (*http.Response, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fileParamName, filepath.Base(path))

	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, ref.addr+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	for _, headerRaw := range headerRaws {
		req.Header.Set(headerRaw.key, headerRaw.value)
	}

	req.WithContext(ctx)

	resp, err := ref.cl.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ref *MultiPartHttpClient) GetFile(ctx context.Context, path string, headerRaws ...HeaderRaw) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, ref.addr+path, nil)
	if err != nil {
		return nil, err
	}

	for _, headerRaw := range headerRaws {
		req.Header.Set(headerRaw.key, headerRaw.value)
	}

	req.WithContext(ctx)

	resp, err := ref.cl.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
