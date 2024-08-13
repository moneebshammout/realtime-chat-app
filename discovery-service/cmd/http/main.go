package main

import (
	"context"
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
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"

	"discovery-service/pkg/utils"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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

func buildGrpcServer() (*grpc.Server, func(), error) {
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

func buildProxyServer() (*runtime.ServeMux, func(), error) {
	// Create the gRPC-Gateway mux
	gatewayMux := runtime.NewServeMux(
		runtime.WithMarshalerOption("application/protobuf", &runtime.ProtoMarshaller{}),
		runtime.WithMarshalerOption("application/json", &runtime.JSONPb{}),
		runtime.WithMetadata(func(_ context.Context, req *http.Request) metadata.MD {
			return metadata.New(map[string]string{
				"x-auth-signature": req.Header.Get("x-auth-signature"),
			})
		}),
	)

	// Register the gRPC services with the gRPC-Gateway mux
	ctx, cancel := context.WithCancel(context.Background())

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := discoveryGRPCGen.RegisterDiscoveryHandlerFromEndpoint(ctx, gatewayMux, fmt.Sprintf(":%s", config.Env.Port), opts)
	if err != nil {
		return nil, nil, err
	}

	return gatewayMux, func() {
		// clean up3
		cancel()
	}, nil
}

func run() (func(), error) {
	// Connect to zookeeper
	go zookeeper.Connect(config.Env.ZooHosts)

	server, cleanup, err := buildGrpcServer()
	if err != nil {
		return nil, err
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	appName := config.Env.App
	go func() {
		port := config.Env.Port
		listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
		if err != nil {
			logger.Panicf("cannot create gRPC listener: %s", err)
		}

		logger.Infof("%s gRPC server running on :%s\n", appName, port)

		if err := server.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			logger.Errorf("Error starting gRPC server: %v\n", err)
		}
	}()

	// Start the HTTP server (with the gRPC-gateway)
	proxyMux, proxyCleanup, err := buildProxyServer()
	if err != nil {
		return nil, err
	}

	go func() {
		port := config.Env.GatewayPort

		logger.Infof("%s HTTP server (with gRPC-gateway) running on :%s\n", appName, port)

		if err := http.ListenAndServe(fmt.Sprintf(":%s", port), proxyMux); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Error starting HTTP server: %v\n", err)
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
		proxyCleanup()
		zookeeper.Close()
	}, nil
}
