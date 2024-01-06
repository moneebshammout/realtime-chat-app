package proxy

import (
	"net/http/httputil"
	"net/url"

	"github.com/labstack/echo/v4"
)

func Proxy(echo *echo.Echo, paths []string, url *url.URL) {
	proxy := httputil.NewSingleHostReverseProxy(url)
	for _, path := range paths {
		echo.Any(path+"/*", reverseProxyHandler(path, proxy))
	}
}