/*
 * Copyright The Titan Project Contributors.
 */
package echo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestType(t *testing.T) {
	r := EchoRemote{}
	typ, _ := r.Type()
	assert.Equal(t, "echo", typ)
}

func TestToURL(t *testing.T) {
	r := EchoRemote{}
	u, props, _ := r.ToURL(map[string]interface{}{"a": "b"})
	assert.Equal(t, "echo://echo", u)
	assert.Len(t, props, 1)
	assert.Equal(t, "b", props["a"])
}

func TestFromURL(t *testing.T) {
	r := EchoRemote{}
	res, _ := r.FromURL("echo://echo", map[string]string{"a": "b"})
	assert.Len(t, res, 2)
	assert.Equal(t, "b", res["a"])
	assert.Equal(t, "echo://echo", res["url"])
}

func TestGetParameters(t *testing.T) {
	r := EchoRemote{}
	res, _ := r.GetParameters(map[string]interface{}{"a": "b"})
	assert.Len(t, res, 1)
	assert.Equal(t, "b", res["a"])
}
