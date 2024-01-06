package middleware

import (
	"slices"

	echojwt "github.com/labstack/echo-jwt/v4"

	"github.com/labstack/echo/v4"

	"gateway/pkg/types"
)

func AuthMiddleware(authConfig types.AuthConfig) echo.MiddlewareFunc {
	skipper := func(c echo.Context) bool {
		return slices.Contains(authConfig.PublicRoutes, c.Request().URL.Path)
	}

	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(authConfig.SigningKey),
		ContextKey:  "token",
		TokenLookup: authConfig.TokenLookup,
		Skipper:     skipper,
	})
}
