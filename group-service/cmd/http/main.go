package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"group-service/config"
	"group-service/internal/database"
	"group-service/internal/database/migrations"
	"group-service/internal/groups"
	"group-service/internal/middleware"
	"group-service/pkg/utils"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	groupGRPC "group-service/internal/gRPC"
	groupGRPCGen "group-service/internal/gRPC/group-service-grpc-gen"
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

	groups.Router(app.Group("/api"))

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

func buildGrpcServer() (*grpc.Server, func(), error) {
	logger.Info("Building GRPC server")
	serverRegistrar := grpc.NewServer(
		grpc.ChainUnaryInterceptor(groupGRPC.Interceptors()...),
	)

	service:= &groupGRPC.GroupServiceServer{}
	groupGRPCGen.RegisterGroupServiceServer(serverRegistrar, service)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(serverRegistrar, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	return serverRegistrar, func() {
		healthServer.Shutdown()
	}, nil
}

func run() (func(), error) {
	app, cleanup, err := buildServer()
	if err != nil {
		logger.Errorf("Error building server: %v\n", err)
		return nil, err
	}

	grpcApp, cleanupGrpc, err := buildGrpcServer()
	if err != nil {
		logger.Errorf("Error building gRPC server: %v\n", err)
		return nil, err
	}

	database.Connect()
	migrations.Migrate()

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

	// Start the gRPC server in a goroutine
	go func() {
		port := config.Env.GRPCPort
		appName := config.Env.App

		h2s := &http2.Server{}
		handler := h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && r.Header.Get("content-type") == "application/grpc" {
				grpcApp.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("Hello from %s GRPC Server", appName)))
			}
		}), h2s)

		listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
		if err != nil {
			logger.Panicf("cannot create GRPC listener: %s", err)
		}

		logger.Infof("%s GRPC Server ----> running on http://localhost:%s\n", appName, port)

		if err := http.Serve(listener, handler); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Error starting GRPC server: %v\n", err)
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

	grpcApp.GracefulStop()

	// Return a function to close the server and perform cleanup
	return func() {
		cleanup()
		database.Disconnect()
		cleanupGrpc()
	}, err
}
