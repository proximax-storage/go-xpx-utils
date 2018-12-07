package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// IsOkResponse checking
// 1. Does http.Response equal nil and
// 2. Does http.Response has http status code 200
// Use it only for tests
func IsOkResponse(t *testing.T, resp *http.Response) bool {
	return IsValidResponse(t, resp, false, http.StatusOK)
}

// IsOkResponse checking:
// 1. Does http.Response equal nil and
// 2. Does http.Response has http status code 400
// Use it only for tests
func IsBadRequestResponse(t *testing.T, resp *http.Response) bool {
	return IsValidResponse(t, resp, false, http.StatusBadRequest)
}

// IsValidResponse checking does http.Response correspond to conditions passed as arguments
// Use it only for tests
func IsValidResponse(t *testing.T, resp *http.Response, canBeNil bool, requiredHttpCode int) bool {
	if !canBeNil {
		assert.NotNil(t, resp, "response is nil")
	} else if resp == nil {
		return true
	}

	assert.Equal(t, requiredHttpCode, resp.StatusCode)

	return true
}

// ValidateStringers compares expected and actual string by using fmt.Stringer.String()
// Use it only for tests
func ValidateStringers(t *testing.T, expected, actual fmt.Stringer) {
	if expected == nil && actual == nil {
		return
	}

	if expected == nil {
		assert.Nil(t, actual)
	} else {
		assert.NotNil(t, actual)
	}

	assert.Equal(t, expected.String(), actual.String())
}
