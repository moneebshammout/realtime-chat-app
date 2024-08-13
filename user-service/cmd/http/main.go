package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"user-service/config"
	"user-service/internal/auth"
	"user-service/internal/middleware"
	"user-service/pkg/utils"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

var logger = utils.InitLogger()

func init() {
	config.DBConnect()
}

func main() {
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	// Run the server
	cleanup, err := run()

	// Run the cleanup after the server is terminated
	defer cleanup()

	if err != nil {
		logger.Errorf("Error: %v\n", err)
		exitCode = 1
		return
	}
}

func buildServer() (*echo.Echo, func(), error) {
	// Echo instance
	app := echo.New()

	app.HTTPErrorHandler = middleware.ErrorMiddleware
	// Middleware
	app.Use(echoMiddleware.Logger())
	app.Use(echoMiddleware.Recover())
	app.Use(middleware.CorsMiddleware([]string{
		config.Env.GatewayHost,
	}))
	// Routes
	auth.Router(app)

	app.GET("/healthCheck", func(c echo.Context) error {
		return c.JSON(200, "Hello From "+config.Env.App+" -----> running on Port:"+config.Env.Port)
	})

	app.Any("*", func(c echo.Context) error {
		return c.JSON(200, "You arrived no where")
	})
	return app, func() {
		// Cleanup logic (if any)
	}, nil
}

func run() (func(), error) {
	app, cleanup, err := buildServer()
	if err != nil {
		return nil, err
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Start the server in a goroutine
	go func() {
		port := config.Env.Port
		appName := config.Env.App
		logger.Infof("%s----> running on http://localhost:%s\n", appName, port)

		if err := app.Start(fmt.Sprintf(":%s", port)); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Error starting server: %v\n", err)
		}
	}()

	// Handle exit signals and gracefully shut down the server
	<-interrupt
	logger.Warnf("Received interrupt signal. Initiating graceful shutdown...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the Echo instance
	if err := app.Shutdown(ctx); err != nil {
		logger.Errorf("Error during server shutdown: %v\n", err)
	}

	// Return a function to close the server and perform cleanup
	return func() {
		cleanup()
		config.KillDBConnection()
	}, err
}
