/*
 * Copyright The Titan Project Contributors.
 */
package remote

import (
	"context"
	"github.com/hashicorp/go-plugin"
	proto "github.com/titan-data/remote-sdk-go/internal/proto"
	"github.com/titan-data/remote-sdk-go/internal/util"
	"google.golang.org/grpc"
)

type remoteRPCClient struct {
	Client proto.RemoteClient
}

func (r remoteRPCClient) Type() (string, error) {
	req := proto.GetTypeRequest{}
	res, err := r.Client.GetType(context.Background(), &req)
	if err != nil {
		return "", err
	}
	return res.Type, nil
}

func (r remoteRPCClient) FromURL(url string, properties map[string]string) (map[string]interface{}, error) {
	req := proto.FromURLRequest{Url: url, Properties: properties}
	res, err := r.Client.FromURL(context.Background(), &req)
	if err != nil {
		return nil, err
	}
	output, err := util.Struct2Map(res.Remote)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (r remoteRPCClient) ToURL(properties map[string]interface{}) (string, map[string]string, error) {
	s, err := util.Map2Struct(properties)
	if err != nil {
		return "", nil, err
	}
	req := proto.ToURLRequest{Remote: s}
	res, err := r.Client.ToURL(context.Background(), &req)
	if err != nil {
		return "", nil, err
	}
	return res.Url, res.Properties, nil
}

func (r remoteRPCClient) GetParameters(properties map[string]interface{}) (map[string]interface{}, error) {
	p, err := util.Map2Struct(properties)
	if err != nil {
		return nil, err
	}
	req := proto.GetParametersRequest{Remote: p}
	res, err := r.Client.GetParameters(context.Background(), &req)
	if err != nil {
		return nil, err
	}
	return util.Struct2Map(res.Parameters)
}

func (r remoteRPCClient) ValidateRemote(properties map[string]interface{}) error {
	p, err := util.Map2Struct(properties)
	if err != nil {
		return err
	}
	req := proto.ValidateRemoteRequest{Remote: p}
	_, err = r.Client.ValidateRemote(context.Background(), &req)
	return err
}

func (r remoteRPCClient) ValidateParameters(parameters map[string]interface{}) error {
	p, err := util.Map2Struct(parameters)
	if err != nil {
		return err
	}
	req := proto.ValidateParametersRequest{Parameters: p}
	_, err = r.Client.ValidateParameters(context.Background(), &req)
	return err
}

func (r remoteRPCClient) ListCommits(properties map[string]interface{}, parameters map[string]interface{}, tags []Tag) ([]Commit, error) {
	remote, err := util.Map2Struct(properties)
	if err != nil {
		return nil, err
	}
	params, err := util.Map2Struct(parameters)
	if err != nil {
		return nil, err
	}
	rpcTags := make([]*proto.Tag, len(tags))
	for i, t := range tags {
		if t.Value != nil {
			rpcTags[i] = &proto.Tag{
				Key: t.Key,
				Value: &proto.Tag_ValueString{ValueString: *t.Value},
			}
		} else {
			rpcTags[i] = &proto.Tag{
				Key: t.Key,
				Value: &proto.Tag_ValueNull{ValueNull: true},
			}
		}
	}
	input := proto.ListCommitRequest{
		Properties:           remote,
		Paramegers:           params,
		Tags:                 rpcTags,
	}
	res, err := r.Client.ListCommits(context.Background(), &input)
	if err != nil {
		return nil, err
	}
	nativeCommits := make([]Commit, len(res.Commits))
	for i, c := range res.Commits {
		props, err := util.Struct2Map(c.Properties)
		if err != nil {
			return nil, err
		}
		nativeCommits[i] = Commit{
			Id:         c.Id,
			Properties: props,
		}
	}
	return nativeCommits, nil
}

func (r remoteRPCClient) GetCommit(properties map[string]interface{}, parameters map[string]interface{}, commitId string) (*Commit, error) {
	remote, err := util.Map2Struct(properties)
	if err != nil {
		return nil, err
	}
	params, err := util.Map2Struct(parameters)
	if err != nil {
		return nil, err
	}
	input := proto.GetCommitRequest{
		Properties:           remote,
		Parameters:           params,
		CommitId:             commitId,
	}
	res, err := r.Client.GetCommit(context.Background(), &input)
	if err != nil {
		return nil, err
	}
	if res.GetCommitNull() {
		return nil, nil
	} else {
		props, err := util.Struct2Map(res.GetCommitValue().Properties)
		if err != nil {
			return nil, err
		}
		return &Commit { Id: commitId, Properties: props}, nil
	}
}

type remoteRPCServer struct {
	Impl Remote
}

func (r *remoteRPCServer) GetType(context.Context, *proto.GetTypeRequest) (*proto.GetTypeResponse, error) {
	typ, err := r.Impl.Type()
	if err != nil {
		return nil, err
	}
	return &proto.GetTypeResponse{Type: typ}, nil
}

func (r *remoteRPCServer) FromURL(ctx context.Context, req *proto.FromURLRequest) (*proto.FromURLResponse, error) {
	props, err := r.Impl.FromURL(req.Url, req.Properties)
	if err != nil {
		return nil, err
	}
	output, err := util.Map2Struct(props)
	if err != nil {
		return nil, err
	}
	return &proto.FromURLResponse{Remote: output}, nil
}

func (r *remoteRPCServer) ToURL(ctx context.Context, req *proto.ToURLRequest) (*proto.ToURLResponse, error) {
	input, err := util.Struct2Map(req.Remote)
	if err != nil {
		return nil, err
	}
	url, props, err := r.Impl.ToURL(input)
	if err != nil {
		return nil, err
	}
	ret := &proto.ToURLResponse{Url:    url, Properties: props }
	return ret, nil
}

func (r *remoteRPCServer) GetParameters(ctx context.Context, req *proto.GetParametersRequest) (*proto.GetParametersResponse, error) {
	input, err := util.Struct2Map(req.Remote)
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
	return &proto.GetParametersResponse{Parameters: output}, nil
}

func (r *remoteRPCServer) ValidateRemote(ctx context.Context, req *proto.ValidateRemoteRequest) (*proto.ValidateRemoteResponse, error) {
	remote, err := util.Struct2Map(req.Remote)
	if err != nil {
		return nil, err
	}
	err = r.Impl.ValidateRemote(remote)
	return &proto.ValidateRemoteResponse{}, err
}

func (r *remoteRPCServer) ValidateParameters(ctx context.Context, req *proto.ValidateParametersRequest) (*proto.ValidateParametersResponse, error) {
	params, err := util.Struct2Map(req.Parameters)
	if err != nil {
		return nil, err
	}
	err = r.Impl.ValidateParameters(params)
	return &proto.ValidateParametersResponse{}, err
}

func (r *remoteRPCServer) ListCommits(ctx context.Context, req *proto.ListCommitRequest) (*proto.ListCommitResponse, error) {
	remote, err := util.Struct2Map(req.Properties)
	if err != nil {
		return nil, err
	}
	params, err := util.Struct2Map(req.Paramegers)
	if err != nil {
		return nil, err
	}
	nativeTags := make([]Tag, len(req.Tags))
	for i, t := range req.Tags {
		if t.GetValueNull() {
			nativeTags[i] = Tag { Key: t.Key }
		} else {
			val := t.GetValueString()
			nativeTags[i] = Tag { Key: t.Key, Value: &val}
		}
	}
	commits, err := r.Impl.ListCommits(remote, params, nativeTags)
	if err != nil {
		return nil, err
	}

	rpcCommits := make([]*proto.Commit, len(commits))
	for i, c := range commits {
		props, err := util.Map2Struct(c.Properties)
		if err != nil {
			return nil, err
		}
		rpcCommits[i] = &proto.Commit{
			Id:                   c.Id,
			Properties:           props,
		}
	}

	return &proto.ListCommitResponse{ Commits: rpcCommits }, nil
}

func (r *remoteRPCServer) GetCommit(ctx context.Context, req *proto.GetCommitRequest) (*proto.GetCommitResponse, error) {
	remote, err := util.Struct2Map(req.Properties)
	if err != nil {
		return nil, err
	}
	params, err := util.Struct2Map(req.Parameters)
	if err != nil {
		return nil, err
	}
	commit, err := r.Impl.GetCommit(remote, params, req.CommitId)
	if err != nil {
		return nil, err
	}
	if commit == nil {
		return &proto.GetCommitResponse{
			Commit:               &proto.GetCommitResponse_CommitNull{CommitNull: true},
		}, nil
	} else {
		s, err := util.Map2Struct(commit.Properties)
		if err != nil {
			return nil, err
		}
		rpcCommit := proto.Commit{
			Id:                   commit.Id,
			Properties:           s,
		}
		return &proto.GetCommitResponse{
			Commit:               &proto.GetCommitResponse_CommitValue{CommitValue: &rpcCommit},
		}, nil
	}
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
