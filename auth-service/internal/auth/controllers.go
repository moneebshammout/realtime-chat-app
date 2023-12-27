package auth

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

// Handler
func register(c echo.Context) error {
	return c.String(http.StatusOK, "Register")
}

func login(c echo.Context) error {
	return c.String(http.StatusOK, "Login")
}
