package groups

import (
	"group-service/internal/middleware"

	"github.com/labstack/echo/v4"
)

func Router(app *echo.Group) {
	logger.Infof("Adding Group Routes")
	groups := app.Group("/groups")

	groups.GET("/:id", getGroup, middleware.ValidationMiddleware(&GroupGetSerizliser{}))
	groups.POST("", create, middleware.ValidationMiddleware(&GroupCreateSerizliser{}))
}
