package queues

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	appConfig "chat-service/config/app"
	queueConfig "chat-service/config/queues"
	"chat-service/internal/websocket"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynq/x/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SpawnWorkersServer(hub *websocket.Hub) {
	logger.Infof("Starting%s:Queue server", appConfig.Env.App)
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	// Run the server
	cleanup, err := run(hub)

	// Run the cleanup after the server is terminated
	defer cleanup()

	if err != nil {
		logger.Error(err)
		exitCode = 1
		return
	}
}

func buildServer(hub *websocket.Hub) (*asynq.Server, *asynq.ServeMux, func()) {
	logger.Infof("Building%s:Queue server", appConfig.Env.App)

	app := asynq.NewServer(
		asynq.RedisClientOpt{Addr: queueConfig.Env.RedisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				queueConfig.Env.MessageQueue: 6,
			},
		},
	)

	mux := asynq.NewServeMux()
	MessageQueueMux(mux, hub)

	return app, mux, func() {}
}

func run(hub *websocket.Hub) (func(), error) {
	app, mux, cleanup := buildServer(hub)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Start metrics server
	go func() {
		reg := prometheus.NewPedanticRegistry()

		inspector := asynq.NewInspector(asynq.RedisClientOpt{
			Addr: queueConfig.Env.RedisAddr,
		})

		reg.MustRegister(
			metrics.NewQueueMetricsCollector(inspector),
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
			collectors.NewGoCollector(),
		)

		http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
		logger.Infof("Metrics available at http://localhost:%s/metrics", queueConfig.Env.MetricsPort)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", queueConfig.Env.MetricsPort), nil))
	}()

	// Start the main Asynq server
	go func() {
		port := queueConfig.Env.Port
		appName := fmt.Sprintf("%s:Queue server", appConfig.Env.App)
		logger.Infof("%s is running on http://localhost:%s\n", appName, port)
		if err := app.Run(mux); err != nil {
			logger.Errorf("Error starting server: %v\n", err)
			return
		}
	}()

	// Handle exit signals and gracefully shut down the server
	<-interrupt
	logger.Info("Received interrupt signal. Initiating graceful shutdown...")

	// Attempt to gracefully shut down the Asynq instance
	app.Stop()
	app.Shutdown()

	// Return a function to close the server and perform cleanup
	return func() {
		cleanup()
	}, nil
}
