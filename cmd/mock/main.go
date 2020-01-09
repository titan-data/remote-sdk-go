package mock

import (
	"github.com/titan-data/remote-sdk-go/remote"
	"net/url"
)

type MockRemote struct {
}

func (m MockRemote) Type() (string, error) {
	return "mock", nil
}

func (m MockRemote) FromURL(url *url.URL, additionalProperties map[string]string) (map[string]interface{}, error) {
	ret := map[string]interface{}{}
	for k, v := range additionalProperties {
		ret[k] = v
	}
	return ret, nil
}

func (m MockRemote) ToURL(properties map[string]interface{}) (string, map[string]string, error) {
	ret := map[string]string{}
	for k, v := range properties {
		ret[k] = v.(string)
	}
	return "mock://mock", ret, nil
}

func (m MockRemote) GetParameters(remoteProperties map[string]interface{}) (map[string]interface{}, error) {
	return remoteProperties, nil
}

func main() {
	remote.Register(MockRemote{})
	remote.Serve("mock")
}

