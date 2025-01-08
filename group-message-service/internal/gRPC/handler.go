package gRPC

import (
	"context"
	"encoding/json"
	"fmt"

	"group-message-service/config"
	grpcClients "group-message-service/internal/gRPC/clients"
	"group-message-service/pkg/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"group-message-service/internal/clients"
	groupMessageGen "group-message-service/internal/gRPC/group-message-service-grpc-gen"
)

var logger = utils.GetLogger()

type GroupGroupMessageServiceServer struct {
	groupMessageGen.UnimplementedGroupMessageServiceServer
}

type Message struct {
	Message   string
	GroupId   string
	SenderId  string
	CreatedAt int64
}

type Receiver struct {
	Connection   string `json:"connection"`
	Server       string `json:"server"`
	MessageQueue string `json:"messageQueue"`
}

func (s GroupGroupMessageServiceServer) Send(ctx context.Context, req *groupMessageGen.SendRequest) (*groupMessageGen.SendResponse, error) {
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

	groupServiceClient, err := grpcClients.NewGroupServiceClient(config.Env.GroupServiceUrl)
	defer groupServiceClient.Disconnect()
	if err != nil {
		logger.Errorf("Error creating group service client: %v\n", err)
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error creating group service client: %s", err.Error()),
		)
	}

	response, err := groupServiceClient.GetGroupUsers(message.GroupId)
	if err != nil {
		logger.Errorf("Error getting group users: %v\n", err)
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error getting group users: %s", err.Error()),
		)
	}

	for _, userId := range response.UserIds {

		websocketManagerClient, err := grpcClients.NewWebsocketManagerClient(config.Env.WebsocketManagerUrl)
		defer websocketManagerClient.Disconnect()
		if err != nil {
			logger.Errorf("Error creating websocket manager client: %v\n", err)
			return nil, status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error creating websocket manager client: %s", err.Error()),
			)
		}

		response, err := websocketManagerClient.Discover(userId)
		if err != nil {
			// User is not online we send the message to relay service
			// failure in websocket manager also we save the message for late delivery
			logger.Infof("Error discovering Receiver: %v\n err: %v", userId, err)

			jobData := map[string]interface{}{
				"message":    message.Message,
				"senderId":   message.SenderId,
				"receiverId": userId,
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
				"receiverId": userId,
				"createdAt":  message.CreatedAt,
			}

			err = enqueueMessage(jobData, receiver.MessageQueue)
			if err != nil {
				return nil, err
			}
		}
	}

	return &groupMessageGen.SendResponse{
		Status:  "OK",
		Message: "Group Message handeled successfully",
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
