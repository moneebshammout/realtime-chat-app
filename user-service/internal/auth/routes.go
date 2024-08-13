package auth

import (
	"user-service/internal/middleware"

	"github.com/labstack/echo/v4"
)

func Router(app *echo.Echo) {
	logger.Infof("Adding Auth Routes")
	group := app.Group("/auth")

	group.POST("/login", login, middleware.ValidationMiddleware(&LoginSerializer{}))
	group.POST("/register", register, middleware.ValidationMiddleware(&RegisterSerializer{}))
}
