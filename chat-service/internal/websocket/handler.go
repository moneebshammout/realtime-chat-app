package websocket

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.QueryParam("userId")
		if(userId == ""){
			return c.JSON(http.StatusBadRequest, "userId is required")
		}

		conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
		if err != nil {
			logger.Errorf("Error upgrading connection: %v\n", err)
			return err
			// return c.JSON(http.StatusInternalServerError, err.Error())
		}
		client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), userId: userId}
		client.hub.register <- client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.writePump()
		go client.readPump()
		return nil
	}
}
