package groups

import (
	"group-service/internal/middleware"

	"github.com/labstack/echo/v4"
)

func Router(app *echo.Echo) {
	logger.Infof("Adding Group Routes")
	group := app.Group("/group")

	group.POST("/", create, middleware.ValidationMiddleware(&GroupSerizliser{}))
}
