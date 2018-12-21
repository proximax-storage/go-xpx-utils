package net

import (
	"github.com/proximax-storage/proximax-utils-go/mock"
	"github.com/proximax-storage/proximax-utils-go/tests"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
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

	statusCode, err := cl.Get(testContext, "/testGet", inputDTO)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)

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

	statusCode, err := cl.Post(testContext, "/testPost", outputDTO, inputDTO)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)

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

	statusCode, err := cl.Put(testContext, "/testPut", outputDTO, inputDTO)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)

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

	statusCode, err := cl.Delete(testContext, "/testDelete", outputDTO, inputDTO)
	assert.Equal(t, http.StatusOK, statusCode)

	tests.ValidateStringers(t, test1Obj, inputDTO)
}
