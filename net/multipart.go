package net

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type multiPartHttpClient struct {
	addr string
	cl   *http.Client
}

func newMultiPartHttpClient(addr string) (*multiPartHttpClient, error) {
	if len(addr) == 0 {
		return nil, errors.New("address should not be blank")
	}

	return &multiPartHttpClient{addr: addr, cl: &http.Client{}}, nil
}

func (ref *multiPartHttpClient) postFile(ctx context.Context, path, fileParamName, filePath string, options ...RequestOption) (*http.Response, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	r, w := io.Pipe()
	defer r.Close()

	writer := multipart.NewWriter(w)
	go func() {
		defer func() {
			w.Close()
			writer.Close()
		}()

		select {
		case <-ctx.Done():
			return
		default:
		}

		part, err := writer.CreateFormFile(fileParamName, fi.Name())
		if err != nil {
			return
		}

		_, err = io.Copy(part, file) // Could be a io.CopyBuffer for defined chunk size
	}()

	req, err := http.NewRequest(http.MethodPost, ref.addr+path, r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	for _, option := range options {
		option(req)
	}

	req.WithContext(ctx)

	return ref.cl.Do(req)
}

func (ref *multiPartHttpClient) getFile(ctx context.Context, path string, options ...RequestOption) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, ref.addr+path, nil)
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		option(req)
	}

	req.WithContext(ctx)

	return ref.cl.Do(req)
}
