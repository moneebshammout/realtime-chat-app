package auth

import (
	"github.com/labstack/echo/v4"
)

func Router(app *echo.Echo) {
	group := app.Group("/auth")

	group.POST("/login", login)
	group.POST("/register", register)
}
