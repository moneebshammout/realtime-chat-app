package gRPC

import (
	"context"
	"fmt"

	"websocket-manager/config"
	"websocket-manager/internal/clients"
	managerGen "websocket-manager/internal/gRPC/websocket-manager-grpc-gen"

	"websocket-manager/pkg/utils"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var logger = utils.InitLogger()

type WebsocketManagerServer struct {
	managerGen.UnimplementedWebsocketManagerServer
}

func (s WebsocketManagerServer) Register(ctx context.Context, req *managerGen.RegisterRequest) (*managerGen.RegisterResponse, error) {
	redisClient := clients.NewRedisClient(config.Env.RedisUrl)
	defer redisClient.Close()

	err := redisClient.Set(req.UserId, req.Data)
	if err != nil {
		logger.Errorf("Error registering Connection %s with data %s: %s", req.UserId, req.Data, err.Error())
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error registering Connection %s with data %s: %s", req.UserId, req.Data, err.Error()),
		)
	}

	return &managerGen.RegisterResponse{
		Status:  "OK",
		Message: "Connection registered successfully",
	}, nil
}

func (s WebsocketManagerServer) Discover(ctx context.Context, req *managerGen.DiscoverRequest) (*managerGen.DiscoverResponse, error) {
	redisClient := clients.NewRedisClient(config.Env.RedisUrl)
	defer redisClient.Close()

	data, err := redisClient.Get(req.UserId)
	if err != nil {
		if err == redis.Nil {
			logger.Errorf("Connection %s not found", req.UserId)
			return nil, status.Errorf(
				codes.NotFound,
				fmt.Sprintf("Connection %s not found", req.UserId),
			)
		}

		logger.Errorf("Error discovering Connection %s: %s", req.UserId, err.Error())
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error discovering Connection %s: %s", req.UserId, err.Error()),
		)
	}

	return &managerGen.DiscoverResponse{
		Data:    data,
		Status:  "OK",
		Message: "Connection found successfully",
	}, nil
}

func (s WebsocketManagerServer) Unregister(ctx context.Context, req *managerGen.UnregisterRequest) (*managerGen.UnregisterResponse, error) {
	redisClient := clients.NewRedisClient(config.Env.RedisUrl)
	defer redisClient.Close()

	err := redisClient.Del(req.UserId)
	if err != nil {
		logger.Errorf("Error unregistering Connection %s: %s", req.UserId, err.Error())
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error unregistering Connection %s: %s", req.UserId, err.Error()),
		)
	}

	return &managerGen.UnregisterResponse{
		Status:  "OK",
		Message: "Connection unregistered successfully",
	}, nil
}
