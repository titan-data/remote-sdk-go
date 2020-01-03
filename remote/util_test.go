package remote

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func registerDefaultRemote() *MockRemote {
	Clear()
	r := new(MockRemote)
	r.On("Type").Return("mock")
	r.On("FromURL", mock.Anything, mock.Anything).Return(map[string]interface{}{}, nil)
	Register(r)
	return r
}

func TestProvider(t *testing.T) {
	r := registerDefaultRemote()
	provider, _, _, _, _ := ParseURL("mock", map[string]string{})
	assert.Equal(t, "mock", provider)
	r.AssertExpectations(t)
}

func TestBadProvider(t *testing.T) {
	_ = registerDefaultRemote()
	_, _, _, _, err := ParseURL("notmock", map[string]string{})
	assert.NotNil(t, err)
}

func TestBadURL(t *testing.T) {
	_ = registerDefaultRemote()
	_, _, _, _, err := ParseURL("a\000b", map[string]string{})
	assert.NotNil(t, err)
}

func TestProviderScheme(t *testing.T) {
	r := registerDefaultRemote()
	provider, _, _, _, _ := ParseURL("mock://foo", map[string]string{})
	assert.Equal(t, "mock", provider)
	r.AssertExpectations(t)
}

func TestFragment(t *testing.T) {
	r := registerDefaultRemote()
	_, _, _, commit, _ := ParseURL("mock://foo#commit", map[string]string{})
	assert.Equal(t, "commit", commit)
	r.AssertExpectations(t)
}

func TestNoFragment(t *testing.T) {
	r := registerDefaultRemote()
	_, _, _, commit, _ := ParseURL("mock://foo", map[string]string{})
	assert.Empty(t, commit)
	r.AssertExpectations(t)
}

func TestQueryParams(t *testing.T) {
	r := registerDefaultRemote()
	_, _, tags, _, _ := ParseURL("mock://foo?tag=one&tag=two=three", map[string]string{})
	assert.Len(t, tags, 2)
	assert.Equal(t, "one", tags[0])
	assert.Equal(t, "two=three", tags[1])
	r.AssertExpectations(t)
}

func TestBadQueryParams(t *testing.T) {
	_ = registerDefaultRemote()
	_, _, _, _, err := ParseURL("mock://foo?nottag=one", map[string]string{})
	assert.NotNil(t, err)
}

func TestEmptyQueryParams(t *testing.T) {
	r := registerDefaultRemote()
	_, _, tags, _, _ := ParseURL("mock://foo", map[string]string{})
	assert.Empty(t, tags)
	r.AssertExpectations(t)
}

func TestProperties(t *testing.T) {
	Clear()
	r := new(MockRemote)
	r.On("Type").Return("mock")
	r.On("FromURL", mock.Anything, mock.Anything).Return(map[string]interface{}{"a": "b"}, nil)
	Register(r)

	_, props, _, _, _ := ParseURL("mock://foo", map[string]string{})
	assert.Len(t, props, 1)
	assert.Equal(t, "b", props["a"])
	r.AssertExpectations(t)
}

func TestBadProperties(t *testing.T) {
	Clear()
	r := new(MockRemote)
	r.On("Type").Return("mock")
	r.On("FromURL", mock.Anything, mock.Anything).Return(map[string]interface{}{}, errors.New("error"))
	Register(r)

	_, _, _, _, err := ParseURL("mock://foo", map[string]string{})
	assert.NotNil(t, err)
	r.AssertExpectations(t)
}

func TestProviderArguments(t *testing.T) {
	Clear()
	r := new(MockRemote)
	r.On("Type").Return("mock")
	r.On("FromURL", mock.Anything, mock.Anything).Return(map[string]interface{}{}, nil)
	Register(r)

	_, _, _, _, _ = ParseURL("mock://user:pass@host:80/path", map[string]string{"a": "b"})
	assert.Len(t, r.props, 1)
	assert.Equal(t, "b", r.props["a"])
	assert.Equal(t, "mock", r.u.Scheme)
	assert.Equal(t, "host", r.u.Hostname())
	assert.Equal(t, "80", r.u.Port())
	assert.Equal(t, "/path", r.u.Path)
	assert.Equal(t, "user", r.u.User.Username())
	pass, set := r.u.User.Password()
	assert.Equal(t, "pass", pass)
	assert.True(t, set)
	r.AssertExpectations(t)
}
