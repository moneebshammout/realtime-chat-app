package middleware

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func ValidationMiddleware(vStruct interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			v := validator.New()
			
			if err := c.Bind(vStruct); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			}

			// Validate the request body
			if err := v.Struct(vStruct); err != nil {
				// Validation failed
				return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
			}

			// Append the validated data to the context
			c.Set("validatedData", vStruct)

			// Continue to the next handler
			return next(c)
		}
	}
}
