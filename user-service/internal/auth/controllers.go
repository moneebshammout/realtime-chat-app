package auth

import (
	"net/http"
	"time"
	"user-service/config"
	"user-service/pkg/types"
	"user-service/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)


func register(c echo.Context) error {
	return c.String(http.StatusOK, "Register")
}

func login(c echo.Context) error {
	//check password logic and email logic

	// Set custom claims
	claims := &types.JwtCustomClaims{
		Name: "Jon Snow",
		Admin: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(config.Env.JWTAccessExpirayMinutes))),
		},
	}

	accessToken,err := utils.GenerateJWT(config.Env.JWTAccessSecret,*claims)
	if err != nil {
		return err
	}

	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(config.Env.JWTRefreshExpirayHours))),
	}

	refreshToken,err := utils.GenerateJWT(config.Env.JWTRefreshSecret,*claims)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access_token": accessToken,
		"refresh_token" : refreshToken,
	})
}
