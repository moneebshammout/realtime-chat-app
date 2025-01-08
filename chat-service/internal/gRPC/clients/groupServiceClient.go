package clients

import (
	"context"

	appConfig "chat-service/config/app"
	groupMessageGRPCGen "chat-service/internal/gRPC/group-message-service-grpc-gen"
	"chat-service/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type GroupMessageServiceClient struct {
	conn   *grpc.ClientConn
	path   string
	client groupMessageGRPCGen.GroupMessageServiceClient
}

func NewGroupMessageServiceClient(path string) (*GroupMessageServiceClient, error) {
	logger.Infof("GroupMessageServiceClient Connecting to %s\n", path)
	conn, err := grpc.Dial(path, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GroupMessageServiceClient{
		conn:   conn,
		path:   path,
		client: groupMessageGRPCGen.NewGroupMessageServiceClient(conn),
	}, nil
}

func (dc *GroupMessageServiceClient) Disconnect() error {
	err := dc.conn.Close()
	if err != nil {
		logger.Errorf("Error closing GroupMessageServiceClient connection: %v\n", err)
		return err
	}

	logger.Infof("GroupMessageServiceClient Disconnected from %s\n", dc.path)
	return nil
}

func (dc *GroupMessageServiceClient) Send(data string) error {
	payload := &groupMessageGRPCGen.SendRequest{Message: data}
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

	logger.Infof("GroupMessageServiceClient Response: %+v", response)

	return nil
}
