/*
 * Copyright The Titan Project Contributors.
 */
package remote

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegister(t *testing.T) {
	Clear()
	r := new(MockRemote)
	r.On("Type").Return("mock")
	Register(r)

	res := Get("mock")
	assert.Equal(t, "mock", res.Type())
	r.AssertExpectations(t)
}
