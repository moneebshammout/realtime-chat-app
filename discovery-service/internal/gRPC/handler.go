package discovery

import (
	"context"
	"discovery-service/internal/zookeeper"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type DiscoveryServiceServer struct {
	UnimplementedDiscoveryServer
}

func (s DiscoveryServiceServer) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	name := req.Service.Name
	url := req.Service.Url
	err := zookeeper.Register(req.Service.Path, []byte(url))

	if err != nil {
		return nil, status.Errorf(
			status.Code(err),
			fmt.Sprintf("Error registering service %s with url %s: %s", name, url, err.Error()),
		)
	}

	return &RegisterResponse{
		Status:  "OK",
		Message: fmt.Sprintf("Service %s registered with url %s", name, url),
	}, nil
}

func (s DiscoveryServiceServer) Discover(ctx context.Context, req *DiscoverRequest) (*DiscoverResponse, error) {
	//get data from zookeeper
	return &DiscoverResponse{}, nil
}

func (s DiscoveryServiceServer) UnaryServerInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	//interceptor logic
	return handler(ctx, req)
}
