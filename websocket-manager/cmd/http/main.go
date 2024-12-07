package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"websocket-manager/config"

	managerGRPC "websocket-manager/internal/gRPC"
	manageryGRPCGen "websocket-manager/internal/gRPC/websocket-manager-grpc-gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"websocket-manager/pkg/utils"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var logger = utils.InitLogger()

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

func buildServer() (*grpc.Server, func(), error) {
	logger.Infof("buildServer %s\n", config.Env.App)
	serverRegistrar := grpc.NewServer(
		grpc.ChainUnaryInterceptor(managerGRPC.Interceptors()...),
	)

	service := &managerGRPC.WebsocketManagerServer{}
	manageryGRPCGen.RegisterWebsocketManagerServer(serverRegistrar, service)

	// Register gRPC health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(serverRegistrar, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	return serverRegistrar, func() {
		logger.Infof("Shutting down server %s\n", config.Env.App)
		healthServer.Shutdown()
	}, nil
}

func run() (func(), error) {
	server, cleanup, err := buildServer()
	if err != nil {
		logger.Errorf("Error: %v\n", err)
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
			logger.Panicf("cannot create listener: %s", err)
		}

		logger.Infof("%s----> running on http://localhost:%s\n", appName, port)

		if err := http.Serve(listener, handler); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Error starting server: %v\n", err)
		}
	}()

	// Handle exit signals and gracefully shut down the servers
	<-interrupt
	logger.Warnf("Received interrupt signal. Initiating graceful shutdown...")

	// Attempt to gracefully shut down the gRPC server
	server.GracefulStop()

	// Return a function to close the server and perform cleanup
	return func() {
		cleanup()
	}, nil
}
