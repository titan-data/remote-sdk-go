/*
 * Copyright The Titan Project Contributors.
 */
package remote

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPluginType(t *testing.T) {
	Clear()
	r, err := Load("echo", ".")
	if err != nil {
		t.Fatal(err)
	}
	typ, _ := r.Type()
	assert.Equal(t, "echo", typ)
}
