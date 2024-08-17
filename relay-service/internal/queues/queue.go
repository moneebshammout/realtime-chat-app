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

func done(message string, data any, job *asynq.Task) error {
	result := map[string]interface{}{
		"job_id":  job.ResultWriter().TaskID(),
		"status":  "success",
		"message": message,
		"result":  data,
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		logger.Errorf("Error marshalling result: %v", err)
		return err
	}

	_, err = job.ResultWriter().Write(jsonData)

	return err
}
