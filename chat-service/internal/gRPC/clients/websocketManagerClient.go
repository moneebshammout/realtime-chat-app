package clients

import (
	"context"

	appConfig "chat-service/config/app"

	websocketManagerGRPCGen "chat-service/internal/gRPC/websocket-manager-grpc-gen"
	"chat-service/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

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

func (dc *WebsocketManagerClient) Register(userId string, data string) error {
	payload := &websocketManagerGRPCGen.RegisterRequest{UserId: userId, Data: data}
	payloadJson, err := proto.Marshal(payload)
	if err != nil {
		return err
	}

	signature := utils.GenerateHmacSignature(payloadJson, appConfig.Env.SignatureKey)
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

	logger.Infoln("WebsocketManagerClient Response:", response)

	return nil
}

func (dc *WebsocketManagerClient) UnRegister(userId string) error {
	payload := &websocketManagerGRPCGen.UnregisterRequest{UserId: userId}
	payloadJson, err := proto.Marshal(payload)
	if err != nil {
		return err
	}

	signature := utils.GenerateHmacSignature(payloadJson, appConfig.Env.SignatureKey)
	md := metadata.Pairs(
		"x-auth-signature", signature,
	)

	// Create a new context with the metadata
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Make a request with the context containing the metadata
	response, err := dc.client.Unregister(ctx, payload)
	if err != nil {
		return err
	}

	logger.Infoln("WebsocketManagerClient Response:", response)

	return nil
}
