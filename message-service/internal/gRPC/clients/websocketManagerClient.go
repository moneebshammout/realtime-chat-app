package clients

import (
	"context"

	"message-service/config"
	websocketManagerGRPCGen "message-service/internal/gRPC/websocket-manager-grpc-gen"
	"message-service/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

var logger = utils.GetLogger()

type WebsocketManagerClient struct {
	conn   *grpc.ClientConn
	path   string
	client websocketManagerGRPCGen.WebsocketManagerClient
}

func NewWebsocketManagerClient(path string) (*WebsocketManagerClient, error) {
	logger.Infof("WebsocketManagerClient Connecting to %s\n", path)
	conn, err := grpc.Dial(path, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &WebsocketManagerClient{
		conn:   conn,
		path:   path,
		client: websocketManagerGRPCGen.NewWebsocketManagerClient(conn),
	}, nil
}

func (dc *WebsocketManagerClient) Disconnect() error {
	err := dc.conn.Close()
	if err != nil {
		logger.Errorf("Error closing WebsocketManagerClient connection: %v\n", err)
		return err
	}

	logger.Infof("WebsocketManagerClient Disconnected from %s\n", dc.path)
	return nil
}

func (dc *WebsocketManagerClient) Discover(userId string) (*websocketManagerGRPCGen.DiscoverResponse, error) {
	payload := &websocketManagerGRPCGen.DiscoverRequest{UserId: userId}
	payloadJson, err := proto.Marshal(payload)
	if err != nil {
		return nil, err
	}

	signature := utils.GenerateHmacSignature(payloadJson, config.Env.SignatureKey)
	md := metadata.Pairs(
		"x-auth-signature", signature,
	)

	// Create a new context with the metadata
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Make a request with the context containing the metadata
	response, err := dc.client.Discover(ctx, payload)

	logger.Infoln("WebsocketManagerClient Discover Response:", response)

	return response, err
}
