/*
 * Copyright The Titan Project Contributors.
 */
package remote

import (
	"context"
	"github.com/hashicorp/go-plugin"
	"github.com/titan-data/remote-sdk-go/internal/proto"
	"google.golang.org/grpc"
	"net/url"
)

type remoteRPCClient struct {
	Client proto.RemoteClient
}

func (r remoteRPCClient) Type() (string, error) {
	resp, err := r.Client.Type(context.Background(), &proto.Empty{})
	return resp.Type, err
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

func (r *remoteRPCServer) Type(ctx context.Context, req *proto.Empty) (*proto.RemoteType, error) {
	typ, err := r.Impl.Type()
	if err != nil {
		return nil, err
	}
	return &proto.RemoteType{Type: typ}, nil
}

type remotePlugin struct {
	plugin.Plugin
	Impl Remote
}

func (p *remotePlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterRemoteServer(s, &remoteRPCServer{Impl: p.Impl})
	return nil
}

func (remotePlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &remoteRPCClient{Client: proto.NewRemoteClient(c)}, nil
}
