package queues

import (
	"context"
	"fmt"

	queuesConfig "relay-service/config/queues"

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
}

func consumer(ctx context.Context, job *asynq.Task) error {
	logger.Infof("New %s Job received id: %s", queuesConfig.Env.RelayQueue, job.ResultWriter().TaskID())
	var data payload
	if err := getData(job, &data); err != nil {
		return err
	}

	job.ResultWriter().Write([]byte(fmt.Sprintf("Message Stored in DB for User: %s", data.ReceiverId)))

	return nil
}

func RelayQueueMux(mux *asynq.ServeMux) {
	mux.HandleFunc(RelayQueue.Name, consumer)
}
