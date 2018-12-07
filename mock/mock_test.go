package mock

import (
	"github.com/proximax-storage/proximax-utils-go/tests"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

const (
	testPath = "/test"
	testBody = "TestBody"
)

func TestNewMock(t *testing.T) {
	mockServer := NewMock(0)
	defer mockServer.Close()

	resp, err := http.Get(mockServer.GetServerURL() + testPath)

	assert.Nilf(t, err, "http.Get returned error: %s", err)
	tests.IsValidResponse(t, resp, false, http.StatusNotFound)
}

func TestNewMockWithRoute(t *testing.T) {
	mockServer := NewMockWithRoute(&Router{
		Path: testPath,
	})

	defer mockServer.Close()

	resp, err := http.Get(mockServer.GetServerURL() + testPath)

	assert.Nilf(t, err, "http.Get returned error: %s", err)
	tests.IsOkResponse(t, resp)
}

func TestMock_AddHandler(t *testing.T) {
	mockServer := NewMock(0)
	defer mockServer.Close()

	mockServer.AddHandler(testPath, func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(200)
	})

	resp, err := http.Get(mockServer.GetServerURL() + testPath)

	assert.Nilf(t, err, "http.Get returned error: %s", err)
	tests.IsOkResponse(t, resp)
}

func TestMock_AddRouter(t *testing.T) {
	mockServer := NewMock(0)
	defer mockServer.Close()

	mockServer.AddRouter(&Router{
		Path:     testPath,
		RespBody: testBody,
	})

	resp, err := http.Get(mockServer.GetServerURL() + testPath)
	assert.Nilf(t, err, "http.Get returned error: %s", err)

	respBody, err := ioutil.ReadAll(resp.Body)
	assert.Nilf(t, err, "ioutil.ReadAll returned error: %s", err)

	tests.IsOkResponse(t, resp)
	assert.Equal(t, testBody, string(respBody))
}
