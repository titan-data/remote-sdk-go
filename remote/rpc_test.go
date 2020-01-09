/*
 * Copyright The Titan Project Contributors.
 */
package remote

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func getEcho(t *testing.T) Remote {
	Clear()
	r, err := Load("echo", ".")
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestPluginType(t *testing.T) {
	e := getEcho(t)
	ret, err := e.Type()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "echo", ret)
}

func TestToURL(t *testing.T) {
	e := getEcho(t)
	url, props, err := e.ToURL(map[string]interface{}{"a": "b"})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "echo://echo", url)
	assert.Len(t, props, 1)
	assert.Equal(t, "b", props["a"])
}
