/*
 * Copyright The Titan Project Contributors.
 */
package remote

import (
	"github.com/hashicorp/go-plugin"
	"net/rpc"
	"net/url"
)

type remoteRPCClient struct {
	Client *rpc.Client
}

func (r remoteRPCClient) Type() (string, error) {
	var resp string
	err := r.Client.Call("Plugin.Type", new(interface{}), &resp)
	return resp, err
}

func (r remoteRPCClient) FromURL(url *url.URL, additionalProperties map[string]string) (map[string]interface{}, error) {
	panic("implement me")
}

func (r remoteRPCClient) ToURL(properties map[string]interface{}) (string, map[string]string, error) {
	panic("implement me")
}

func (r remoteRPCClient) GetParameters(remoteProperties map[string]interface{}) (map[string]interface{}, error) {
	panic("implement me")
}

type remoteRPCServer struct {
	Impl Remote
}

func (r *remoteRPCServer) Type(args interface{}, resp *string) error {
	var err error
	*resp, err = r.Impl.Type()
	return err
}

type remotePlugin struct {
	Impl Remote
}

func (p *remotePlugin) Server(broker *plugin.MuxBroker) (interface{}, error) {
	return &remoteRPCServer{Impl: p.Impl}, nil
}

func (remotePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &remoteRPCClient{Client: c}, nil
}
