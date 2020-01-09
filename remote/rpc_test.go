/*
 * Copyright The Titan Project Contributors.
 */
package remote

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPluginType(t *testing.T) {
	r, client, err := Load("echo", ".")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Kill()
	typ, _ := r.Type()
	assert.Equal(t, "echo", typ)
}
