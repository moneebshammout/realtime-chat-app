package queues

import (
	"encoding/json"
	"fmt"

	"relay-service/pkg/utils"

	"github.com/hibiken/asynq"
)

var logger = utils.GetLogger()

type Queue struct {
	Name     string
	Consumer any
}

func getData(job *asynq.Task, mapper interface{}) error {
	if err := json.Unmarshal(job.Payload(), &mapper); err != nil {
		err := fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		logger.Error(err)
		return err
	}
	return nil
}
