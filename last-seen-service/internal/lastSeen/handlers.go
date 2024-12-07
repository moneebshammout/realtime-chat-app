package lastSeen

import (
	"fmt"
	"time"

	"last-seen-service/config"
	"last-seen-service/internal/clients"
	"last-seen-service/pkg/utils"

	"github.com/labstack/echo/v4"
)

var logger = utils.GetLogger()

func createLastSeen(c echo.Context) error {
	data := c.Get("validatedData").(*CreateSerizliser)
	redisClient := clients.NewRedisClient(config.Env.RedisUrl)
	defer redisClient.Close()
	err := redisClient.SetX(fmt.Sprintf("%v", data.UserId), data, time.Hour*1)
	if err != nil {
		logger.Errorf("Failed to set last seen: %s", err.Error())
		return c.JSON(400, "Failed to set last seen")
	}

	return c.JSON(200, data)
}

func getLastSeen(c echo.Context) error {
	id := c.Get("validatedData").(*IDParam).ID
	redisClient := clients.NewRedisClient(config.Env.RedisUrl)
	defer redisClient.Close()

	lastSeen, err := redisClient.Get(id)
	if err != nil {
		logger.Errorf("Failed to get last seen: %s", err)
		return c.JSON(404, "Failed to get last seen")
	}

	if lastSeen == "" {
		logger.Errorf("Last seen not found for user: %s", id)
		return c.JSON(404, "Last seen not found")
	}

	return c.JSON(200, lastSeen)
}
