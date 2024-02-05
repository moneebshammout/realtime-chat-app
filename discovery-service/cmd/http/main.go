package main

import (
	"discovery-service/config"
	discovery "discovery-service/internal/gRPC"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"google.golang.org/grpc"
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

func buildServer() (*grpc.Server, func(), error) {
	serverRegistrar := grpc.NewServer()
	service := &discovery.DiscoveryServiceServer{}
	discovery.RegisterDiscoveryServer(serverRegistrar, service)

	return serverRegistrar, func() {
		// Cleanup logic (if any)
	}, nil
}

func run() (func(), error) {
	server, cleanup, err := buildServer()
	if err != nil {
		return nil, err
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Start the server in a goroutine
	go func() {
		port := config.Env.Port
		appName := config.Env.App
		listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
		if err != nil {
			panic(fmt.Sprintf("cannot create listener: %s", err))
		}

		fmt.Printf("%s----> running on http://localhost:%s\n", appName, port)

		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	// Handle exit signals and gracefully shut down the server
	select {
	case <-interrupt:
		fmt.Println("Received interrupt signal. Initiating graceful shutdown...")

		// Create a context with timeout for graceful shutdown
		// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// defer cancel()

		// Attempt to gracefully shut down the grpc instance
		server.GracefulStop()

	}

	// Return a function to close the server and perform cleanup
	return func() {
		cleanup()
	}, nil
}
