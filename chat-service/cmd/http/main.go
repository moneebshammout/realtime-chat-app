package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"chat-service/config"
	"chat-service/internal/middleware"
	"chat-service/internal/websocket"

	grpcClients "chat-service/internal/gRPC/clients"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"chat-service/pkg/utils"
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

func buildServer() (*echo.Echo, *websocket.Hub, func(), error) {
	logger.Info("Building server")
	// Echo instance
	app := echo.New()
	hub := websocket.NewHub()

	// Middleware
	app.HTTPErrorHandler = middleware.ErrorMiddleware
	app.Use(echoMiddleware.Logger())
	app.Use(echoMiddleware.Recover())
	// TODO: uncomment
	// app.Use(middleware.CorsMiddleware([]string{
	// 	config.Env.GatewayHost,
	// }))

	// Routes

	websocket.Router(app, hub)
	app.Any("*", func(c echo.Context) error {
		return c.JSON(200, "You arrived no where")
	})
	return app, hub, func() {
		// Cleanup logic (if any)
	}, nil
}

func registerServer() {
	logger.Info("Registering server")
	time.Sleep(5 * time.Second)
	for {
		discoveryClient, err := grpcClients.NewDiscoveryClient(config.Env.DiscoveryServiceUrl)
		defer discoveryClient.Disconnect()

		if err != nil {
			logger.Errorf("Error creating discovery client: %v\n retrying in 5 second...\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		data := map[string]string{
			"address":  fmt.Sprintf("%s:%s", config.Env.Host, config.Env.Port),
			"location": "amman/jo",
			"name":     config.Env.App,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			logger.Errorf("Error marshalling data: %v\n retrying in 5 seconds...\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		err = discoveryClient.Register("/chats", string(jsonData))
		if err != nil {
			logger.Errorf("Error registering server: %v\n retrying in 5 seconds...\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Println("Server registered successfully")
		break
	}
}

func run() (func(), error) {
	app, hub, cleanup, err := buildServer()
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

		// go registerServer()
		if err := app.Start(fmt.Sprintf(":%s", port)); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Error starting server: %v\n", err)
			return
		}
	}()

	// start websocket HUB in a goroutine
	go hub.Run()

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
	}, nil
}
