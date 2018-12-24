package net

import (
	"bytes"
	ctx "context"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
)

type HttpStatusCode int

type RestClient struct {
	mPartClient *MultiPartHttpClient
	addr        string
	cl          *http.Client
}

func NewRestClient(addr string) (*RestClient, error) {
	if len(addr) == 0 {
		return nil, errors.New("address should not be blank")
	}

	mPartClient, err := NewMultiPartHttpClientClient(addr)
	if err != nil {
		return nil, err
	}

	cl := &RestClient{
		mPartClient: mPartClient,
		addr:        addr,
		cl:          &http.Client{},
	}

	return cl, nil
}

func (ref *RestClient) Get(ctx ctx.Context, path string, inputDTO interface{}, headerRaws ...HeaderRow) (*http.Response, error) {
	return ref.doRequest(ctx, http.MethodGet, path, nil, inputDTO, headerRaws...)
}

func (ref *RestClient) Post(ctx ctx.Context, path string, outputDTO, inputDTO interface{}, headerRaws ...HeaderRow) (*http.Response, error) {
	return ref.doRequest(ctx, http.MethodPost, path, outputDTO, inputDTO, headerRaws...)
}

func (ref *RestClient) Put(ctx ctx.Context, path string, outputDTO, inputDTO interface{}, headerRaws ...HeaderRow) (*http.Response, error) {
	return ref.doRequest(ctx, http.MethodPut, path, outputDTO, inputDTO, headerRaws...)
}

func (ref *RestClient) Delete(ctx ctx.Context, path string, outputDTO, inputDTO interface{}, headerRaws ...HeaderRow) (*http.Response, error) {
	return ref.doRequest(ctx, http.MethodDelete, path, outputDTO, inputDTO, headerRaws...)
}

func (ref *RestClient) PostFile(ctx ctx.Context, path string, fileParamName, filePath string, inputDTO interface{}, headerRaws ...HeaderRow) (*http.Response, error) {
	resp, err := ref.mPartClient.PostFile(ctx, path, fileParamName, filePath, headerRaws...)
	if err != nil {
		return nil, err
	}

	return resp, convertRespToJson(resp.Body, inputDTO)
}

func (ref *RestClient) GetFile(ctx ctx.Context, path string, headerRaws ...HeaderRow) (*http.Response, error) {
	return ref.mPartClient.GetFile(ctx, path, headerRaws...)
}

func (ref *RestClient) doRequest(ctx ctx.Context, method, path string, outputDTO, inputDTO interface{}, headerRaws ...HeaderRow) (*http.Response, error) {
	var body io.Reader

	if outputDTO != nil {
		buf, err := json.Marshal(outputDTO)
		if err != nil {
			return nil, err
		}

		body = bytes.NewReader(buf)
	}

	req, err := http.NewRequest(method, ref.addr+path, body)
	if err != nil {
		return nil, err
	}

	if method != http.MethodGet {
		req.Header.Set("Content-Type", "application/json")
	}

	for _, headerRaw := range headerRaws {
		req.Header.Set(headerRaw.Key, headerRaw.Value)
	}

	req.WithContext(ctx)

	resp, err := ref.cl.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, convertRespToJson(resp.Body, inputDTO)
}

func convertRespToJson(respBody io.Reader, inputDTO interface{}) error {
	if inputDTO != nil {
		if buf, err := ioutil.ReadAll(respBody); err != nil {
			return err
		} else {
			return json.Unmarshal(buf, inputDTO)
		}
	}

	return nil
}
