package gRPC

import (
	"context"
	"encoding/json"
	"fmt"

	"message-service/config"
	grpcClients "message-service/internal/gRPC/clients"
	messageGen "message-service/internal/gRPC/message-service-grpc-gen"
	"message-service/pkg/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var logger = utils.GetLogger()

type MessageServiceServer struct {
	messageGen.UnimplementedDiscoveryServer
}

type Message struct {
	Message    string
	SenderId   string
	ReceiverId string
}

type Receiver struct {
	Connection string `json:"connection"`
	Server     string `json:"server"`
}

func (s MessageServiceServer) Send(ctx context.Context, req *messageGen.SendRequest) (*messageGen.SendResponse, error) {
	logger.Infof("Received message: %s", req.Message)

	message := Message{}
	err := json.Unmarshal([]byte(req.Message), &message)
	if err != nil {
		logger.Errorf("Error unmarshalling message: %s", err.Error())
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error unmarshalling message: %s", err.Error()),
		)
	}

	websocketManagerClient, err := grpcClients.NewWebsocketManagerClient(config.Env.WebsocketManagerUrl)
	defer websocketManagerClient.Disconnect()
	if err != nil {
		logger.Errorf("Error creating websocket manager client: %v\n", err)
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error creating websocket manager client: %s", err.Error()),
		)
	}

	response, err := websocketManagerClient.Discover(message.ReceiverId)
	if err != nil {

		logger.Infof("Error discovering Receiver: %v\n err: %v", message.ReceiverId, err)
		if grpcStatus, ok := status.FromError(err); ok {
			code := grpcStatus.Code()
			logger.Infof("gRPC status code: %v\n", code)

			if code == codes.NotFound {
				logger.Infof("User is not online sending to relay service: %v\n", message.ReceiverId)
				//:TODO enqueue in relay service queue
				return &messageGen.SendResponse{
					Status:  "OK",
					Message: "Message handeled successfully",
				}, nil
			}

		} else {
			return nil, status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error discovering Receiver: %v\n err: %s", message.ReceiverId, err.Error()),
			)
		}

	}

	receiver := Receiver{}
	err = json.Unmarshal([]byte(response.Data), &receiver)
	if err != nil {
		logger.Errorf("Error unmarshalling receiver: %s", err.Error())
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error unmarshalling receiver: %s", err.Error()),
		)
	}

	logger.Infof("Sending message to: %s", receiver.Server)
	logger.Infof("Connection: %s", receiver.Connection)
	//enqueue in receiver server message queue

	return &messageGen.SendResponse{
		Status:  "OK",
		Message: "Message handeled successfully",
	}, nil
}
