package gRPC

import (
	"context"
	"fmt"

	"discovery-service/internal/zookeeper"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DiscoveryServiceServer struct {
	UnimplementedDiscoveryServer
}

func (s DiscoveryServiceServer) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	url := req.Url
	err := zookeeper.Register(req.Path, url)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error registering service with url %s: %s", url, err.Error()),
		)
	}

	return &RegisterResponse{
		Status:  "OK",
		Message: fmt.Sprintf("Service registered with url %s", url),
	}, nil
}

func (s DiscoveryServiceServer) Discover(ctx context.Context, req *DiscoverRequest) (*DiscoverResponse, error) {
	// get data from zookeeper
	data, err := zookeeper.Discover(req.Path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error discovering services: %s", err.Error()),
		)
	}

	if len(data) == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("No services found for path %s", req.Path),
		)
	}

	// put the data in array of service struct

	return &DiscoverResponse{
		Urls:    data,
		Status:  "OK",
		Message: "Services discovered",
	}, nil
}

// func (s DiscoveryServiceServer) UnaryServerInterceptor(ctx context.Context,
// 	req interface{},
// 	info *grpc.UnaryServerInfo,
// 	handler grpc.UnaryHandler,
// ) (interface{}, error) {
// 	// interceptor logic
// 	return handler(ctx, req)
// }
