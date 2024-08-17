package queues

import (
	"context"
	"encoding/json"
	"fmt"

	queuesConfig "relay-service/config/queues"

	"relay-service/internal/clients"
	"relay-service/internal/database/models"

	"github.com/hibiken/asynq"
)

var RelayQueue Queue

func init() {
	RelayQueue = Queue{
		Name:     queuesConfig.Env.RelayQueue,
		Consumer: consumer,
	}
}

type payload struct {
	Message    string `json:"message"`
	SenderId   string `json:"senderId"`
	ReceiverId string `json:"receiverId"`
	CreatedAt  int64  `json:"createdAt"`
}

func consumer(ctx context.Context, job *asynq.Task) error {
	logger.Infof("New %s Job received id: %s", queuesConfig.Env.RelayQueue, job.ResultWriter().TaskID())
	var data payload
	if err := getData(job, &data); err != nil {
		return err
	}

	messageDAO := models.MessageDAO()
	messageObject := map[string]interface{}{
		"message":     data.Message,
		"sender_id":   data.SenderId,
		"receiver_id": data.ReceiverId,
		"created_at":  data.CreatedAt,
	}

	result, err := messageDAO.Create(messageObject)
	if err != nil {
		logger.Errorf("Failed to create message: %v", err)
		return err
	}

	//notifiy user about message
	notificationQueue := clients.NewQueueClient(queuesConfig.Env.RedisAddr, queuesConfig.Env.NotificationQueue)
	defer notificationQueue.Close()
	notification := map[string]interface{}{
		"message_id":  result["id"],
		"receiver_id": data.ReceiverId,
		"created_at":  data.CreatedAt,
	}
	notificationJson, _ := json.Marshal(notification)
	notificationQueue.Enqueue(notificationJson)

	return done(fmt.Sprintf("Message %s Stored in DB", result["id"]), result, job)
}

func RelayQueueMux(mux *asynq.ServeMux) {
	mux.HandleFunc(RelayQueue.Name, consumer)
}
