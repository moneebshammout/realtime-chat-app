package middleware

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
)

func CorsMiddleware(allowedHosts []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the request host
			host := c.Request().Host
			if !slices.Contains(allowedHosts, host) {
				// If not, deny the request
				return c.JSON(http.StatusForbidden, fmt.Sprintf("Access denied. Host:%s not allowed.", host))
			}

			// If the host is allowed, continue to the next handler
			return next(c)
		}
	}
}
