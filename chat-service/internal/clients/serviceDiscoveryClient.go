package clients

import (
	"context"
	"fmt"

	"chat-service/config"
	discoveryGRPCGen "chat-service/internal/gRPC/discovery-grpc-gen"
	"chat-service/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type DiscoveryClient struct {
	conn   *grpc.ClientConn
	path   string
	client discoveryGRPCGen.DiscoveryClient
}

func NewDiscoveryClient(path string) (*DiscoveryClient, error) {
	fmt.Printf("Connecting to %s\n", path)
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
	return dc.conn.Close()
}

func (dc *DiscoveryClient) Register(path string, data string) error {
	payload:=&discoveryGRPCGen.RegisterRequest{Path: path, Data: data}
	payloadJson,err := proto.Marshal(payload)
	if err != nil {
		return err
	}

	signature:=utils.GenerateHmacSignature(payloadJson, config.Env.SignatureKey)
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

	fmt.Println("Response:", response)
	

	return nil
}
