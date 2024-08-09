package clients

import (
	"context"

	"chat-service/config"
	discoveryGRPCGen "chat-service/internal/gRPC/discovery-grpc-gen"
	"chat-service/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

var logger = utils.GetLogger()

type DiscoveryClient struct {
	conn   *grpc.ClientConn
	path   string
	client discoveryGRPCGen.DiscoveryClient
}

func NewDiscoveryClient(path string) (*DiscoveryClient, error) {
	logger.Infof("DiscoveryClient Connecting to %s\n", path)
	conn, err := grpc.Dial(path, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &DiscoveryClient{
		conn:   conn,
		path:   path,
		client: discoveryGRPCGen.NewDiscoveryClient(conn),
	}, nil
}

func (dc *DiscoveryClient) Disconnect() error {
	err := dc.conn.Close()
	if err != nil {
		logger.Errorf("Error closing DiscoveryClient connection: %v\n", err)
		return err
	}

	logger.Infof("DiscoveryClient Disconnected from %s\n", dc.path)
	return nil
}

func (dc *DiscoveryClient) Register(path string, data string) error {
	payload := &discoveryGRPCGen.RegisterRequest{Path: path, Data: data}
	payloadJson, err := proto.Marshal(payload)
	if err != nil {
		return err
	}

	signature := utils.GenerateHmacSignature(payloadJson, config.Env.SignatureKey)
	md := metadata.Pairs(
		"x-auth-signature", signature,
	)

	// Create a new context with the metadata
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Make a request with the context containing the metadata
	response, err := dc.client.Register(ctx, payload)
	if err != nil {
		return err
	}

	logger.Infof("DiscoveryClient Response: %+v", response)

	return nil
}
