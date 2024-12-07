package clients

import (
	"context"

	appConfig "chat-service/config/app"
	messageGRPCGen "chat-service/internal/gRPC/message-service-grpc-gen"
	"chat-service/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type MessageServiceClient struct {
	conn   *grpc.ClientConn
	path   string
	client messageGRPCGen.MessageServiceClient
}

func NewMessageServiceClient(path string) (*MessageServiceClient, error) {
	logger.Infof("MessageServiceClient Connecting to %s\n", path)
	conn, err := grpc.Dial(path, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &MessageServiceClient{
		conn:   conn,
		path:   path,
		client: messageGRPCGen.NewMessageServiceClient(conn),
	}, nil
}

func (dc *MessageServiceClient) Disconnect() error {
	err := dc.conn.Close()
	if err != nil {
		logger.Errorf("Error closing MessageServiceClient connection: %v\n", err)
		return err
	}

	logger.Infof("MessageServiceClient Disconnected from %s\n", dc.path)
	return nil
}

func (dc *MessageServiceClient) Send(data string) error {
	payload := &messageGRPCGen.SendRequest{Message: data}
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
	response, err := dc.client.Send(ctx, payload)
	if err != nil {
		return err
	}

	logger.Infof("MessageServiceClient Response: %+v", response)

	return nil
}
