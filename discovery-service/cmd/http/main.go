package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"discovery-service/config"
	discoverGRPC "discovery-service/internal/gRPC"
	discoveryGRPCGen "discovery-service/internal/gRPC/discovery-grpc-gen"

	"discovery-service/internal/zookeeper"

    "google.golang.org/grpc"
    "google.golang.org/grpc/health"
    "google.golang.org/grpc/health/grpc_health_v1"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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
	serverRegistrar := grpc.NewServer(
		grpc.ChainUnaryInterceptor(discoverGRPC.Interceptors()...),
	)

	service := &discoverGRPC.DiscoveryServiceServer{}
	discoveryGRPCGen.RegisterDiscoveryServer(serverRegistrar, service)

	// Register gRPC health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(serverRegistrar, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	return serverRegistrar, func() {
		healthServer.Shutdown()
	}, nil
}

func run() (func(), error) {
	// Connect to zookeeper
	zookeeper.Connect(config.Env.ZooHosts)

	server, cleanup, err := buildServer()
	if err != nil {
		return nil, err
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Start the multiplexed server in a goroutine
	go func() {
		port := config.Env.Port
		appName := config.Env.App

		// Create a handler that routes to the gRPC server and the HTTP router
		h2s := &http2.Server{}
		handler := h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && r.Header.Get("content-type") == "application/grpc" {
				server.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("Hello from %s", appName)))
			}
		}), h2s)

		listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
		if err != nil {
			panic(fmt.Sprintf("cannot create listener: %s", err))
		}

		fmt.Printf("%s----> running on http://localhost:%s\n", appName, port)

		if err := http.Serve(listener, handler); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	// Handle exit signals and gracefully shut down the servers
	<-interrupt
	fmt.Println("Received interrupt signal. Initiating graceful shutdown...")

	// Attempt to gracefully shut down the gRPC server
	server.GracefulStop()

	// Return a function to close the server and perform cleanup
	return func() {
		cleanup()
		zookeeper.Close()
	}, nil
}
