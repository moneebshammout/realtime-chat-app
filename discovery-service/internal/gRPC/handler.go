package gRPC

import (
	"context"
	"fmt"

	"discovery-service/internal/zookeeper"
	"discovery-service/pkg/utils"

	discoveryGen "discovery-service/internal/gRPC/discovery-grpc-gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var logger = utils.GetLogger()

type DiscoveryServiceServer struct {
	discoveryGen.UnimplementedDiscoveryServer
}

func (s DiscoveryServiceServer) Register(ctx context.Context, req *discoveryGen.RegisterRequest) (*discoveryGen.RegisterResponse, error) {
	data := req.Data
	err := zookeeper.Register(req.Path, data)
	if err != nil {
		logger.Errorf("Error registering service with data %s: %s", data, err.Error())
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error registering service with data %s: %s", data, err.Error()),
		)
	}

	return &discoveryGen.RegisterResponse{
		Status:  "OK",
		Message: fmt.Sprintf("Service registered with data %s", data),
	}, nil
}

func (s DiscoveryServiceServer) Discover(ctx context.Context, req *discoveryGen.DiscoverRequest) (*discoveryGen.DiscoverResponse, error) {
	// get data from zookeeper
	data, err := zookeeper.Discover(req.Path)
	if err != nil {
		logger.Errorf("Error discovering services: %s", err.Error())
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error discovering services: %s", err.Error()),
		)
	}

	if len(data) == 0 {
		logger.Errorf("No services found for path %s", req.Path)
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("No services found for path %s", req.Path),
		)
	}

	// put the data in array of service struct

	return &discoveryGen.DiscoverResponse{
		Nodes:   data,
		Status:  "OK",
		Message: "Services discovered",
	}, nil
}
