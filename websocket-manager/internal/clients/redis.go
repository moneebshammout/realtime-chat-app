package clients

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(url string) *RedisClient {
	opts, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opts)

	return &RedisClient{
		client: client,
	}
}

func (c *RedisClient) Close() error {
	return c.client.Close()
}

func (c *RedisClient) Set(key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = c.client.Set(context.TODO(), key, data, 0).Err()
	if err != nil {
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
