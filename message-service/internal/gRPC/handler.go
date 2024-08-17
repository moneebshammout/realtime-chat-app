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

	"message-service/internal/clients"
)

var logger = utils.GetLogger()

type MessageServiceServer struct {
	messageGen.UnimplementedDiscoveryServer
}

type Message struct {
	Message    string
	SenderId   string
	ReceiverId string
	CreatedAt  int64
}

type Receiver struct {
	Connection   string `json:"connection"`
	Server       string `json:"server"`
	MessageQueue string `json:"messageQueue"`
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
		// User is not online we send the message to relay service
		// failure in websocket manager also we save the message for late delivery
		logger.Infof("Error discovering Receiver: %v\n err: %v", message.ReceiverId, err)

		//:TODO enqueue in relay service queue
		jobData := map[string]interface{}{
			"message":    message.Message,
			"senderId":   message.SenderId,
			"receiverId": message.ReceiverId,
			"createdAt":  message.CreatedAt,
		}

		logger.Infof("Enqueuing message: %v", jobData)

		err = enqueueMessage(jobData, config.Env.RelayQueue)
		if err != nil {
			return nil, err
		}

	} else {
		// User is online
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

		jobData := map[string]interface{}{
			"message":    message.Message,
			"senderId":   message.SenderId,
			"receiverId": message.ReceiverId,
			"createdAt":  message.CreatedAt,
		}

		err = enqueueMessage(jobData, receiver.MessageQueue)
		if err != nil {
			return nil, err
		}
	}

	return &messageGen.SendResponse{
		Status:  "OK",
		Message: "Message handeled successfully",
	}, nil
}

func enqueueMessage(data map[string]interface{}, queueName string) error {
	messageQueue := clients.NewQueueClient(config.Env.RedisUrl, queueName)
	defer messageQueue.Close()

	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("Error marshalling data: %v\n", err)
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error marshalling data: %s", err.Error()),
		)
	}

	_, err = messageQueue.Enqueue(jsonData)
	if err != nil {
		logger.Errorf("Error enqueing to %s %v\n", queueName, err)
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error enqueing to %s %s", queueName, err.Error()),
		)
	}

	return nil
}
