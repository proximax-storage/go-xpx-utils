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
	mPartClient *multiPartHttpClient
	addr        string
	cl          *http.Client
}

func NewRestClient(addr string) (*RestClient, error) {
	if len(addr) == 0 {
		return nil, errors.New("address should not be blank")
	}

	mPartClient, err := newMultiPartHttpClient(addr)
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

func (ref *RestClient) Get(ctx ctx.Context, path string, inputDTO interface{}, options ...RequestOption) (*http.Response, error) {
	return ref.doRequest(ctx, http.MethodGet, path, nil, inputDTO, options...)
}

func (ref *RestClient) Post(ctx ctx.Context, path string, outputDTO, inputDTO interface{}, options ...RequestOption) (*http.Response, error) {
	return ref.doRequest(ctx, http.MethodPost, path, outputDTO, inputDTO, options...)
}

func (ref *RestClient) Put(ctx ctx.Context, path string, outputDTO, inputDTO interface{}, options ...RequestOption) (*http.Response, error) {
	return ref.doRequest(ctx, http.MethodPut, path, outputDTO, inputDTO, options...)
}

func (ref *RestClient) Delete(ctx ctx.Context, path string, outputDTO, inputDTO interface{}, options ...RequestOption) (*http.Response, error) {
	return ref.doRequest(ctx, http.MethodDelete, path, outputDTO, inputDTO, options...)
}

func (ref *RestClient) PostFile(ctx ctx.Context, path string, fileParamName, filePath string, inputDTO interface{}, options ...RequestOption) (*http.Response, error) {
	resp, err := ref.mPartClient.postFile(ctx, path, fileParamName, filePath, options...)
	if err != nil {
		return nil, err
	}

	return resp, handleResponse(resp, inputDTO)
}

func (ref *RestClient) GetFile(ctx ctx.Context, path string, options ...RequestOption) (*http.Response, error) {
	resp, err := ref.mPartClient.getFile(ctx, path, options...)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusInternalServerError {
		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return resp, convertError(respBodyBytes)
	}

	return resp, nil
}

func (ref *RestClient) doRequest(ctx ctx.Context, method, path string, outputDTO, inputDTO interface{}, options ...RequestOption) (*http.Response, error) {
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

	for _, option := range options {
		option(req)
	}

	req.WithContext(ctx)

	resp, err := ref.cl.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, handleResponse(resp, inputDTO)
}

func handleResponse(resp *http.Response, dto interface{}) error {
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	status := resp.StatusCode

	switch status {
	case http.StatusOK, http.StatusAccepted:
		return convertRespToJson(respBodyBytes, dto)

	case http.StatusBadRequest, http.StatusInternalServerError:
		return convertError(respBodyBytes)
	}

	return errors.New(string(respBodyBytes))
}

func convertRespToJson(respBody []byte, inputDTO interface{}) error {
	if inputDTO != nil {
		return json.Unmarshal(respBody, inputDTO)
	}

	return nil
}

func convertError(respBody []byte) error {
	identifiableError := &IdentifiableError{}

	if err := json.Unmarshal(respBody, identifiableError); err != nil {
		return err
	}

	return identifiableError
}
