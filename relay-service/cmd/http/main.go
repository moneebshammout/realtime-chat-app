package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	appConfig "relay-service/config/app"
	"relay-service/internal/messages"
	"relay-service/internal/middleware"

	"relay-service/internal/queues"
	"relay-service/pkg/utils"

	"relay-service/internal/database"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

var logger = utils.InitLogger()

func main() {
	logger.Info("Starting server")
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	// Run the server
	cleanup, err := run()

	// Run the cleanup after the server is terminated
	defer cleanup()

	if err != nil {
		utils.GetLogger().Error(err)
		exitCode = 1
		return
	}
}

func buildServer() (*echo.Echo, func(), error) {
	logger.Info("Building server")
	// Echo instance
	app := echo.New()

	// Middleware
	app.HTTPErrorHandler = middleware.ErrorMiddleware
	app.Use(echoMiddleware.Logger())
	app.Use(echoMiddleware.Recover())
	app.Use(middleware.CorsMiddleware([]string{
		appConfig.Env.GatewayHost,
	}))

	// Routes
	messages.Router(app)
	
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

	database.Connect()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Start the server in a goroutine
	go func() {
		port := appConfig.Env.Port
		appName := appConfig.Env.App
		logger.Infof("%s----> running on http://localhost:%s\n", appName, port)

		if err := app.Start(fmt.Sprintf(":%s", port)); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Error starting server: %v\n", err)
			return
		}
	}()

	// run the queues workers server
	go queues.SpawnWorkersServer()

	// Handle exit signals and gracefully shut down the server

	<-interrupt
	logger.Info("Received interrupt signal. Initiating graceful shutdown...")

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
		database.Disconnect()
	}, nil
}

