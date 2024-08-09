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

	"chat-service/internal/clients"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

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
		fmt.Printf("Error: %v\n", err)
		exitCode = 1
		return
	}
}

func buildServer() (*echo.Echo, *websocket.Hub, func(), error) {
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

func registerServer(){
	time.Sleep(5 * time.Second)
	for {
		discoveryClient,err := clients.NewDiscoveryClient(config.Env.DiscoveryServiceUrl)
		defer discoveryClient.Disconnect()

		if err != nil {
			fmt.Printf("Error creating discovery client: %v\n retrying in 5 second...\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		data:= map[string]string{
			"address": fmt.Sprintf("%s:%s",config.Env.Host, config.Env.Port),
			"location": "amman/jo",
			"name": config.Env.App,
		}

		jsonData,err:=json.Marshal(data)
		if err != nil {
			fmt.Printf("Error marshalling data: %v\n retrying in 5 seconds...\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		err = discoveryClient.Register("/chats", string(jsonData))
		if err != nil {
			fmt.Printf("Error registering server: %v\n retrying in 5 seconds...\n", err)
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
		fmt.Printf("%s----> running on http://localhost:%s\n", appName, port)
		
		go registerServer()
		if err := app.Start(fmt.Sprintf(":%s", port)); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
			return
		}
		

		
	}()

	// start websocket HUB in a goroutine
	go hub.Run()

	// Handle exit signals and gracefully shut down the server

	<-interrupt
	fmt.Println("Received interrupt signal. Initiating graceful shutdown...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the Echo instance
	if err := app.Shutdown(ctx); err != nil {
		fmt.Printf("Error during server shutdown: %v\n", err)
	}

	// Return a function to close the server and perform cleanup
	return func() {
		cleanup()
	}, nil
}
