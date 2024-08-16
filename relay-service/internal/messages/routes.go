package messages



import (
	"relay-service/internal/middleware"

	"github.com/labstack/echo/v4"
)

func Router(app *echo.Echo) {
	logger.Infof("Adding Messages Routes")
	group := app.Group("/messages")

	group.GET("/:id", getUserMessages, middleware.ValidationMiddleware(&IGetUserMessages{}))
}
