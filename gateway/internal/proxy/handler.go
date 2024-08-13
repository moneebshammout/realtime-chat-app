package proxy

import (
	"gateway/pkg/utils"
	"net/http/httputil"

	"github.com/labstack/echo/v4"
)

var logger = utils.GetLogger()
func reverseProxyHandler(prefix string, proxy *httputil.ReverseProxy) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Forward the request to the backend service
		logger.Infof("Forwarding request to %s", prefix)
		proxy.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
