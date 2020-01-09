/*
 * Copyright The Titan Project Contributors.
 */
package remote

import (
	"github.com/stretchr/testify/mock"
	"net/url"
)

type MockRemote struct {
	mock.Mock

	u     *url.URL
	props map[string]string
}

func (r *MockRemote) Type() (string, error) {
	args := r.Called()
	return args.String(0), nil
}

func (r *MockRemote) FromURL(url *url.URL, additionalProperties map[string]string) (map[string]interface{}, error) {
	r.u = url
	r.props = additionalProperties

	args := r.Called(url, additionalProperties)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (r *MockRemote) ToURL(properties map[string]interface{}) (string, map[string]string, error) {
	args := r.Called(properties)
	return args.String(0), args.Get(1).(map[string]string), args.Error(2)
}

func (r *MockRemote) GetParameters(remoteProperties map[string]interface{}) (map[string]interface{}, error) {
	args := r.Called(remoteProperties)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}
