package queues

import (
	"context"
	"fmt"

	queuesConfig "chat-service/config/queues"
	"chat-service/internal/websocket"

	"github.com/hibiken/asynq"
)

var MessageQueue Queue

func init() {
	MessageQueue = Queue{
		Name:     queuesConfig.Env.MessageQueue,
		Consumer: consumer,
	}
}

type payload struct {
	Message    string `json:"message"`
	SenderId   string `json:"senderId"`
	ReceiverId string `json:"receiverId"`
}

func consumer(hub *websocket.Hub) func(ctx context.Context, job *asynq.Task) error {
	return func(ctx context.Context, job *asynq.Task) error {
		logger.Infof("New %s Job received id: %s", queuesConfig.Env.MessageQueue, job.ResultWriter().TaskID())
		var data payload
		if err := getData(job, &data); err != nil {
			return err
		}

		hub.Send <- &websocket.SendMessage{
			SenderId:   data.SenderId,
			RecevierId: data.ReceiverId,
			Message:    fmt.Sprintf("from queue %s: %s", data.SenderId, data.Message),
		}
		
		job.ResultWriter().Write([]byte(fmt.Sprintf("Message Sent to receiver: %s", data.ReceiverId)))
	
		return nil
	}
}

func MessageQueueMux(mux *asynq.ServeMux, hub *websocket.Hub) {
	var consumer func(ctx context.Context, job *asynq.Task) error = MessageQueue.Consumer.(func(hub *websocket.Hub) func(ctx context.Context, job *asynq.Task) error)(hub)
	mux.HandleFunc(MessageQueue.Name, consumer)
}
