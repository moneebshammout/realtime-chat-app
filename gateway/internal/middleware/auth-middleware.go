package middleware

import (
	echojwt "github.com/labstack/echo-jwt/v4"

	"github.com/labstack/echo/v4"

	"gateway/pkg/types"
)

func AuthMiddleware(authConfig types.AuthConfig) echo.MiddlewareFunc {
	skipper := func(c echo.Context) bool {
		for _, publicRoute := range authConfig.PublicRoutes {
			if publicRoute == c.Request().URL.Path {
				return true // Skip JWT middleware for public routes
			}
		}
		return false
	}

	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(authConfig.SigningKey),
		ContextKey:  "token",
		TokenLookup: authConfig.TokenLookup,
		Skipper:     skipper,
	})
}
