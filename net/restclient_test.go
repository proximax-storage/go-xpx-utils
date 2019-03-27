package net

import (
	"context"
	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/proximax-storage/go-xpx-utils/tests"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const (
	testRespBody1 = `{"msg":"Hello World"}`
)

var (
	testContext = context.Background()
	test1Obj    = &testOne{
		Msg: "Hello World",
	}
)

type testOne struct {
	Msg string `json:"msg"`
}

func (ref *testOne) String() string {
	return ref.Msg
}

func TestRestClient_Get(t *testing.T) {
	mockServer := mock.NewMockWithRoute(&mock.Router{
		Path:                "/testGet",
		RespHttpCode:        http.StatusOK,
		RespBody:            testRespBody1,
		AcceptedHttpMethods: []string{http.MethodGet},
	})
	defer mockServer.Close()

	cl, err := NewRestClient(mockServer.GetServerURL())
	assert.Nil(t, err)

	inputDTO := &testOne{}

	response, err := cl.Get(testContext, "/testGet", inputDTO)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	tests.ValidateStringers(t, test1Obj, inputDTO)
}

func TestRestClient_Post(t *testing.T) {
	mockServer := mock.NewMockWithRoute(&mock.Router{
		Path:                "/testPost",
		RespHttpCode:        http.StatusOK,
		RespBody:            testRespBody1,
		ReqJsonBodyStruct:   testOne{},
		AcceptedHttpMethods: []string{http.MethodPost},
	})
	defer mockServer.Close()

	cl, err := NewRestClient(mockServer.GetServerURL())
	assert.Nil(t, err)

	outputDTO := &testOne{Msg: "Hello"}
	inputDTO := &testOne{}

	response, err := cl.Post(testContext, "/testPost", outputDTO, inputDTO)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	tests.ValidateStringers(t, test1Obj, inputDTO)
}

func TestRestClient_PostFile(t *testing.T) {
	mockServer := mock.NewMockWithRoute(&mock.Router{
		Path:         "/testPostFile",
		RespHttpCode: http.StatusOK,
		RespBody:     testRespBody1,
		FormParams: []mock.FormParameter{
			{
				"file",
				false,
			},
		},
		AcceptedHttpMethods: []string{http.MethodPost},
	})
	defer mockServer.Close()

	cl, err := NewRestClient(mockServer.GetServerURL())
	assert.Nil(t, err)

	inputDTO := &testOne{}

	response, err := cl.PostFile(testContext, "/testPostFile", "file", "restclient_test.go", inputDTO)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	tests.ValidateStringers(t, test1Obj, inputDTO)
}

func TestRestClient_Put(t *testing.T) {
	mockServer := mock.NewMockWithRoute(&mock.Router{
		Path:                "/testPut",
		RespHttpCode:        http.StatusOK,
		RespBody:            testRespBody1,
		ReqJsonBodyStruct:   testOne{},
		AcceptedHttpMethods: []string{http.MethodPut},
	})
	defer mockServer.Close()

	cl, err := NewRestClient(mockServer.GetServerURL())
	assert.Nil(t, err)

	outputDTO := &testOne{Msg: "Hello"}
	inputDTO := &testOne{}

	response, err := cl.Put(testContext, "/testPut", outputDTO, inputDTO)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	tests.ValidateStringers(t, test1Obj, inputDTO)
}

func TestRestClient_Delete(t *testing.T) {
	mockServer := mock.NewMockWithRoute(&mock.Router{
		Path:                "/testDelete",
		RespHttpCode:        http.StatusOK,
		RespBody:            testRespBody1,
		ReqJsonBodyStruct:   testOne{},
		AcceptedHttpMethods: []string{http.MethodDelete},
	})
	defer mockServer.Close()

	cl, err := NewRestClient(mockServer.GetServerURL())
	assert.Nil(t, err)

	outputDTO := &testOne{Msg: "Hello"}
	inputDTO := &testOne{}

	response, err := cl.Delete(testContext, "/testDelete", outputDTO, inputDTO)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	tests.ValidateStringers(t, test1Obj, inputDTO)
}
