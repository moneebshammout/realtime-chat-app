package websocket

import (
	"github.com/labstack/echo/v4"
)

func Router(app *echo.Echo, hub *Hub) {
	app.GET("/ws", ServeWs(hub))
	app.GET("/", func(c echo.Context) error {
		return c.File("template/home.html")
	})
}
