package mock

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestType(t *testing.T) {
	r := MockRemote{}
	typ, _ := r.Type()
	assert.Equal(t, "mock", typ)
}

func TestToURL(t *testing.T) {
	r := MockRemote{}
	u, props, _ := r.ToURL(map[string]interface{}{"a": "b"})
	assert.Equal(t, "mock://mock", u)
	assert.Len(t, props, 1)
	assert.Equal(t, "b", props["a"])
}

func TestFromURL(t *testing.T) {
	r := MockRemote{}
	u, _ := url.Parse("mock://mock")
	res, _ := r.FromURL(u, map[string]string{"a": "b"})
	assert.Len(t, res, 1)
	assert.Equal(t, "b", res["a"])
}

func TestGetParameters(t *testing.T) {
	r := MockRemote{}
	res, _ := r.GetParameters(map[string]interface{}{"a": "b"})
	assert.Len(t, res, 1)
	assert.Equal(t, "b", res["a"])
}

