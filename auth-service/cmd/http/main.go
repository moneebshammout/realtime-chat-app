package main

import (
	"auth-service/config"
	"auth-service/internal/auth"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func buildServer() (*echo.Echo, func(), error) {
	// Echo instance
	app := echo.New()

	// Middleware
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	//Routes
	auth.Router(app)

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
		var port string = config.Env.Port
		fmt.Printf("Server running on http://localhost:%s\n", port)

		if err := app.Start(fmt.Sprintf(":%s", port)); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	// Handle exit signals and gracefully shut down the server
	select {
	case <-interrupt:
		fmt.Println("Received interrupt signal. Initiating graceful shutdown...")

		// Create a context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Attempt to gracefully shut down the Echo instance
		if err := app.Shutdown(ctx); err != nil {
			fmt.Printf("Error during server shutdown: %v\n", err)
		}
	}

	// Return a function to close the server and perform cleanup
	return func() {
		cleanup()
	}, nil
}
