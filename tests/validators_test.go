package tests

import (
	"github.com/proximax-storage/proximax-utils-go/mock"
	"github.com/proximax-storage/proximax-utils-go/str"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestIsOkResponse(t *testing.T) {
	mockServer := mock.NewMockWithRoute(&mock.Router{
		Path:         "/test1",
		RespHttpCode: 200,
	})

	defer mockServer.Close()

	resp, err := http.Get(mockServer.GetServerURL() + "/test1")

	assert.Nil(t, err)
	IsOkResponse(t, resp)
}

func TestIsBadRequestResponse(t *testing.T) {
	mockServer := mock.NewMockWithRoute(&mock.Router{
		Path:         "/test1",
		RespHttpCode: 400,
	})

	defer mockServer.Close()

	resp, err := http.Get(mockServer.GetServerURL() + "/test1")

	assert.Nil(t, err)
	IsBadRequestResponse(t, resp)
}

func TestIsValidResponse(t *testing.T) {
	mockServer := mock.NewMockWithRoute(&mock.Router{
		Path:         "/test1",
		RespHttpCode: 500,
	})

	defer mockServer.Close()

	resp, err := http.Get(mockServer.GetServerURL() + "/test1")

	assert.Nil(t, err)
	IsValidResponse(t, resp, false, http.StatusInternalServerError)
}

type testStruct struct {
	one string
}

func (t *testStruct) String() string {
	return str.StructToString(
		"testStruct",
		str.NewField("one", str.StringPattern, t.one),
	)
}

func TestValidateStringers(t *testing.T) {
	a := &testStruct{one: "Hello"}
	b := &testStruct{one: "Hello"}

	ValidateStringers(t, a, b)
}
