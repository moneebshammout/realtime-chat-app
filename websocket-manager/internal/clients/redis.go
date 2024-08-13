package clients

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
		"websocket-manager/pkg/utils"
)

var logger = utils.GetLogger()
type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(url string) *RedisClient {
	opts, err := redis.ParseURL(url)
	if err != nil {
		logger.Panicf("Failed to parse Redis URL: %s", err)
	}

	client := redis.NewClient(opts)

	return &RedisClient{
		client: client,
	}
}

func (c *RedisClient) Close() error {
	logger.Info("Closing Redis Client")
	return c.client.Close()
}

func (c *RedisClient) Set(key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		logger.Errorf("Failed to marshal data: %s", err)
		return err
	}

	err = c.client.Set(context.TODO(), key, data, 0).Err()
	if err != nil {
		logger.Errorf("Failed to set data in Redis: %s", err)
		return err
	}

	return nil
}

func (c *RedisClient) Get(key string) (string, error) {
	result, err := c.client.Get(context.TODO(), key).Result()
	if result != "" {
		json.Unmarshal([]byte(result), &result)
	}

	return result, err
}

func (c *RedisClient) Del(key string) error {
	return c.client.Del(context.TODO(), key).Err()
}
