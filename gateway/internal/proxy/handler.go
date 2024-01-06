package proxy

import (
	"net/http/httputil"

	"github.com/labstack/echo/v4"
)

func reverseProxyHandler(prefix string, proxy *httputil.ReverseProxy) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Forward the request to the backend service
		proxy.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
