package websocket

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
		client.hub.register <- client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.writePump()
		go client.readPump()
		return nil
	}
}
