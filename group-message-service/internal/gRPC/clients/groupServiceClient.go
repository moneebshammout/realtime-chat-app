package clients

import (
	"context"

	"group-message-service/config"
	groupServiceGRPCGen "group-message-service/internal/gRPC/group-service-grpc-gen"
	"group-message-service/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)


type GroupServiceClient struct {
	conn   *grpc.ClientConn
	path   string
	client groupServiceGRPCGen.GroupServiceClient
}

func NewGroupServiceClient(path string) (*GroupServiceClient, error) {
	logger.Infof("GroupServiceClient Connecting to %s\n", path)
	conn, err := grpc.Dial(path, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GroupServiceClient{
		conn:   conn,
		path:   path,
		client: groupServiceGRPCGen.NewGroupServiceClient(conn),
	}, nil
}

func (dc *GroupServiceClient) Disconnect() error {
	err := dc.conn.Close()
	if err != nil {
		logger.Errorf("Error closing GroupServiceClient connection: %v\n", err)
		return err
	}

	logger.Infof("GroupServiceClient Disconnected from %s\n", dc.path)
	return nil
}

func (dc *GroupServiceClient) GetGroupUsers(groupId string) (*groupServiceGRPCGen.GetGroupUsersResponse, error) {
	payload := &groupServiceGRPCGen.GetGroupUsersRequest{GroupId: groupId}
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
	response, err := dc.client.GetGroupUsers(ctx, payload)

	logger.Infoln("GroupServiceClient GetGroupUsers Response:", response)

	return response, err
}
