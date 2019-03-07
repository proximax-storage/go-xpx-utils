package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func AssertNil(t *testing.T, obj interface{}, msgAndArgs ...interface{}) {
	if !assert.Nil(t, obj, msgAndArgs...) {
		t.FailNow()
	}
}

func AssertNotNil(t *testing.T, obj interface{}, msgAndArgs ...interface{}) {
	if !assert.NotNil(t, obj, msgAndArgs...) {
		t.FailNow()
	}
}

func AssertEqual(t *testing.T, got, want interface{}, msgAndArgs ...interface{}) {
	if !assert.Equal(t, got, want, msgAndArgs...) {
		t.FailNow()
	}
}

func AssertNotEqual(t *testing.T, got, want interface{}, msgAndArgs ...interface{}) {
	if !assert.NotEqual(t, got, want, msgAndArgs...) {
		t.FailNow()
	}
}

func AssertFalse(t *testing.T, val bool, msgAndArgs ...interface{}) {
	if !assert.False(t, val, msgAndArgs...) {
		t.FailNow()
	}
}

func AssertTrue(t *testing.T, val bool, msgAndArgs ...interface{}) {
	if !assert.True(t, val, msgAndArgs...) {
		t.FailNow()
	}
}

func AssertEmpty(t *testing.T, object interface{}, msgAndArgs ...interface{}) {
	if !assert.Empty(t, object, msgAndArgs...) {
		t.FailNow()
	}
}

func AssertNotEmpty(t *testing.T, object interface{}, msgAndArgs ...interface{}) {
	if !assert.NotEmpty(t, object, msgAndArgs...) {
		t.FailNow()
	}
}
