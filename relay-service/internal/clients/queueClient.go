package clients

import (
	"time"

	"relay-service/pkg/utils"

	"github.com/hibiken/asynq"
)

type QueueClient struct {
	client *asynq.Client
	name   string
}

func NewQueueClient(redisUrl string, name string) *QueueClient {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisUrl})

	return &QueueClient{client: client, name: name}
}

func (c *QueueClient) Close() error {
	return c.client.Close()
}

func (c *QueueClient) Enqueue(payload []byte) (*asynq.TaskInfo, error) {
	task := asynq.NewTask(c.name, payload, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute))
	job, err := c.client.Enqueue(task, asynq.Queue(c.name), asynq.MaxRetry(10), asynq.Timeout(3*time.Minute), asynq.Retention(12*30*24*time.Hour))
	if err != nil {
		utils.GetLogger().Error(err)
	}

	utils.GetLogger().Infof("Enqueued %s job: %s", c.name, job.ID)
	return job, err
}
