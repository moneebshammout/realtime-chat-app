package messages

import (
	"net/http"

	// "relay-service/internal/database/db"
	"relay-service/internal/database/models"
	"relay-service/pkg/utils"

	"github.com/labstack/echo/v4"
	"github.com/scylladb/gocqlx/v3/qb"
)

var logger = utils.GetLogger()

func getUserMessages(c echo.Context) error {
	payload := c.Get("validatedData").(*IGetUserMessages)
	logger.Info(payload)
	messageDAO := models.MessageDAO()
	messages, err := messageDAO.GetList(qb.M{"receiverId =": payload.ReceiverID})
	if err != nil {
		logger.Error(err)
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSONPretty(http.StatusOK, messages, "  ")
}
