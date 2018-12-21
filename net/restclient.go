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

func (ref *RestClient) Get(ctx ctx.Context, path string, inputDTO interface{}, headerRaws ...HeaderRaw) (HttpStatusCode, error) {
	return ref.doRequest(ctx, http.MethodGet, path, nil, inputDTO, headerRaws...)
}

func (ref *RestClient) Post(ctx ctx.Context, path string, outputDTO, inputDTO interface{}, headerRaws ...HeaderRaw) (HttpStatusCode, error) {
	return ref.doRequest(ctx, http.MethodPost, path, outputDTO, inputDTO, headerRaws...)
}

func (ref *RestClient) Put(ctx ctx.Context, path string, outputDTO, inputDTO interface{}, headerRaws ...HeaderRaw) (HttpStatusCode, error) {
	return ref.doRequest(ctx, http.MethodPut, path, outputDTO, inputDTO, headerRaws...)
}

func (ref *RestClient) Delete(ctx ctx.Context, path string, outputDTO, inputDTO interface{}, headerRaws ...HeaderRaw) (HttpStatusCode, error) {
	return ref.doRequest(ctx, http.MethodDelete, path, outputDTO, inputDTO, headerRaws...)
}

func (ref *RestClient) PostFile(ctx ctx.Context, path string, reader io.Reader, inputDTO interface{}, headerRaws ...HeaderRaw) (HttpStatusCode, error) {
	respBody, statusCode, err := ref.mPartClient.PostFile(ctx, path, reader, headerRaws...)
	if err != nil {
		return 0, err
	}

	return statusCode, convertRespToJson(respBody, inputDTO)
}

func (ref *RestClient) GetFile(ctx ctx.Context, path string, headerRaws ...HeaderRaw) (io.Reader, HttpStatusCode, error) {
	return ref.mPartClient.GetFile(ctx, path, headerRaws...)
}

func (ref *RestClient) doRequest(ctx ctx.Context, method, path string, outputDTO, inputDTO interface{}, headerRaws ...HeaderRaw) (HttpStatusCode, error) {
	var (
		buf    []byte
		err    error
		reader io.Reader
	)

	if outputDTO != nil {
		buf, err = json.Marshal(outputDTO)
		if err != nil {
			return 0, err
		}

		reader = bytes.NewReader(buf)
	}

	req, err := http.NewRequest(method, ref.addr+path, reader)
	if err != nil {
		return 0, err
	}

	if method != http.MethodGet {
		req.Header.Set("Content-Type", "application/json")
	}

	for _, headerRaw := range headerRaws {
		req.Header.Set(headerRaw.key, headerRaw.value)
	}

	req.WithContext(ctx)

	resp, err := ref.cl.Do(req)
	if err != nil {
		return 0, err
	}

	return HttpStatusCode(resp.StatusCode), convertRespToJson(resp.Body, inputDTO)
}

func convertRespToJson(respBody io.Reader, inputDTO interface{}) error {
	if inputDTO != nil {
		buf, err := ioutil.ReadAll(respBody)
		if err != nil {
			return err
		}

		return json.Unmarshal(buf, inputDTO)
	}

	return nil
}
