package lastSeen

import (
	"last-seen-service/internal/middleware"

	"github.com/labstack/echo/v4"
)

func Router(app *echo.Group) {
	logger.Infof("Adding Last Seen Routes")
	lastSeens := app.Group("/last-seens")

	lastSeens.GET("/:id", getLastSeen, middleware.ValidationMiddleware(&IDParam{}))
	lastSeens.POST("", createLastSeen, middleware.ValidationMiddleware(&CreateSerizliser{}))
}
