package auth

import (
	"fmt"
	"net/http"
	"time"

	"user-service/config"
	"user-service/internal/database/db"
	"user-service/pkg/types"
	"user-service/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var logger = utils.GetLogger()

func register(c echo.Context) error {
	data := c.Get("validatedData").(*RegisterSerializer)
	hashedPassword, salt, hashError := utils.HashPassword(data.Password)
	if hashError != nil {
		logger.Errorf("could'nt create user %v", hashError.Error())
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("could'nt create user %v", hashError.Error()))
	}

	prisma := config.DB()
	_, createError := prisma.Client.User.CreateOne(
		db.User.Password.Set(hashedPassword),
		db.User.Salt.Set(salt),
		db.User.Email.Set(data.Email),
		db.User.Name.Set(data.Name),
	).Exec(prisma.Context)
	if createError != nil {
		logger.Errorf("could'nt create user %v", createError.Error())
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("could'nt create user %v", createError.Error()))
	}

	return c.JSON(http.StatusCreated, "User Created")
}

func login(c echo.Context) error {
	// check password logic and email logic
	data := c.Get("validatedData").(*LoginSerializer)
	prisma := config.DB()

	user, err := prisma.Client.User.FindFirst(
		db.User.Email.Equals(data.Email),
	).Exec(prisma.Context)
	if err != nil {
		logger.Errorf("could'nt find user %v", err.Error())
		return c.JSON(http.StatusNotFound, err.Error())
	}

	if ok := utils.CheckPassword(data.Password, user.Password, user.Salt); !ok {
		logger.Errorf("could'nt find user %v", err.Error())
		return c.JSON(http.StatusUnauthorized, "Email or Password mismatch!")
	}
	// Set custom claims
	name, _ := user.Name()
	claims := &types.JwtCustomClaims{
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(config.Env.JWTAccessExpirayMinutes))),
		},
	}

	accessToken, err := utils.GenerateJWT(config.Env.JWTAccessSecret, *claims)
	if err != nil {
		logger.Errorf("could'nt generate access token %v", err.Error())
		return err
	}

	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(config.Env.JWTRefreshExpirayHours))),
	}

	refreshToken, err := utils.GenerateJWT(config.Env.JWTRefreshSecret, *claims)
	if err != nil {
		logger.Errorf("could'nt generate refresh token %v", err.Error())
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func refresh(c echo.Context) error {
	claims, err := utils.ValidateJWT(config.Env.JWTRefreshSecret, c.Request().Header.Get("refresh"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid Token")
	}

	accessToken, err := utils.GenerateJWT(config.Env.JWTAccessSecret, *claims)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access_token": accessToken,
	})
}
