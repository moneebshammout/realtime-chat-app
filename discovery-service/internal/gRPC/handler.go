package discovery

import (
	"context"
	"fmt"
)

type DiscoveryServiceServer struct {
	UnimplementedDiscoveryServer
}

func (s DiscoveryServiceServer) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	//get data from zookeeper
	name := req.Service.Name
	url := req.Service.Url

	return &RegisterResponse{
		Status:  "OK",
		Message: fmt.Sprintf("Service %s registered with url %s", name, url),
	}, nil
}

func (s DiscoveryServiceServer) Discover(ctx context.Context, req *DiscoverRequest) (*DiscoverResponse, error) {
	//get data from zookeeper
	return &DiscoverResponse{}, nil
}
