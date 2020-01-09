/*
 * Copyright The Titan Project Contributors.
 */
package remote

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-plugin"
	proto "github.com/titan-data/remote-sdk-go/internal/proto"
	"github.com/titan-data/remote-sdk-go/internal/util"
	"google.golang.org/grpc"
)

type remoteRPCClient struct {
	Client proto.RemoteClient
}

func (r remoteRPCClient) Type() (string, error) {
	res, err := r.Client.Type(context.Background(), &empty.Empty{})
	if err != nil {
		return "", err
	}
	return res.Type, nil
}

func (r remoteRPCClient) FromURL(url string, additionalProperties map[string]string) (map[string]interface{}, error) {
	panic("implement me")
}

func (r remoteRPCClient) ToURL(properties map[string]interface{}) (string, map[string]string, error) {
	s, err := util.Map2Struct(properties)
	if err != nil {
		return "", nil, err
	}
	props := proto.RemoteProperties{Values: s}
	res, err := r.Client.ToURL(context.Background(), &props)
	if err != nil {
		return "", nil, err
	}
	return res.Url, res.Values, nil
}

func (r remoteRPCClient) GetParameters(properties map[string]interface{}) (map[string]interface{}, error) {
	p, err := util.Map2Struct(properties)
	if err != nil {
		return nil, err
	}
	input := proto.RemoteProperties{Values: p}
	res, err := r.Client.GetParameters(context.Background(), &input)
	if err != nil {
		return nil, err
	}
	return util.Struct2Map(res.Values)
}

type remoteRPCServer struct {
	Impl Remote
}

func (r *remoteRPCServer) Type(context.Context, *empty.Empty) (*proto.RemoteType, error) {
	typ, err := r.Impl.Type()
	if err != nil {
		return nil, err
	}
	return &proto.RemoteType{Type: typ}, nil
}

func (r *remoteRPCServer) ToURL(ctx context.Context, req *proto.RemoteProperties) (*proto.ExtendedURL, error) {
	input, err := util.Struct2Map(req.Values)
	if err != nil {
		return nil, err
	}
	url, props, err := r.Impl.ToURL(input)
	if err != nil {
		return nil, err
	}
	ret := &proto.ExtendedURL{
		Url:    url,
		Values: props,
	}
	return ret, nil
}

func (r *remoteRPCServer) GetParameters(ctx context.Context, req *proto.RemoteProperties) (*proto.ParameterProperties, error) {
	input, err := util.Struct2Map(req.Values)
	if err != nil {
		return nil, err
	}
	props, err := r.Impl.GetParameters(input)
	if err != nil {
		return nil, err
	}
	output, err := util.Map2Struct(props)
	if err != nil {
		return nil, err
	}
	return &proto.ParameterProperties{Values: output}, nil
}

type remotePlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl Remote
}

func (p *remotePlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterRemoteServer(s, &remoteRPCServer{Impl: p.Impl})
	return nil
}

func (remotePlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &remoteRPCClient{Client: proto.NewRemoteClient(c)}, nil
}
