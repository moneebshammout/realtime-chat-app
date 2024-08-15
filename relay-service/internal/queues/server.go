package queues

import (
	"fmt"
	// "net/http"
	"os"
	"os/signal"
	"time"

	appConfig "relay-service/config/app"
	queueConfig "relay-service/config/queues"

	"github.com/hibiken/asynq"
	// "github.com/hibiken/asynq/x/metrics"
	// "github.com/hibiken/asynqmon"
	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/collectors"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
	// "github.com/rs/cors"
)

func SpawnWorkersServer() {
	logger.Infof("Starting%s:Queue server", appConfig.Env.App)
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	// Run the server
	cleanup, err := run()

	// Run the cleanup after the server is terminated
	defer cleanup()

	if err != nil {
		logger.Error(err)
		exitCode = 1
		return
	}
}

func buildServer() (*asynq.Server, *asynq.ServeMux, func()) {
	logger.Infof("Building%s:Queue server", appConfig.Env.App)

	app := asynq.NewServer(
		asynq.RedisClientOpt{Addr: queueConfig.Env.RedisAddr},
		asynq.Config{
			Concurrency: 10,

			Queues:          map[string]int{queueConfig.Env.RelayQueue: 6},
			Logger:          logger,
			LogLevel:        asynq.DebugLevel,
			ShutdownTimeout: 5 * time.Second,
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				return 5 * time.Second
			},
			HealthCheckInterval:      5 * time.Second,
			DelayedTaskCheckInterval: 5 * time.Second,
		},
	)

	mux := asynq.NewServeMux()
	RelayQueueMux(mux)

	return app, mux, func() {}
}

// func startMetricsServer() {
// 	queueDashbaord := asynqmon.New(asynqmon.Options{
// 		RootPath: "/queues",
// 		RedisConnOpt: asynq.RedisClientOpt{
// 			Addr:     queueConfig.Env.RedisAddr,
// 			Password: "",
// 			DB:       0,
// 		},
// 		PrometheusAddress: queueConfig.Env.PrometheusAddress,
// 	})

// 	defer queueDashbaord.Close()
// 	c := cors.New(cors.Options{
// 		AllowedMethods: []string{"GET", "POST", "DELETE"},
// 	})
// 	mux := http.NewServeMux()
// 	mux.Handle("/", c.Handler(queueDashbaord))
// 	reg := prometheus.NewPedanticRegistry()

// 	inspector := asynq.NewInspector(asynq.RedisClientOpt{
// 		Addr: queueConfig.Env.RedisAddr,
// 	})

// 	reg.MustRegister(
// 		metrics.NewQueueMetricsCollector(inspector),
// 		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
// 		collectors.NewGoCollector(),
// 	)
// 	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
// 	logger.Infof("Queue Metrics Server available at http://localhost:%s", queueConfig.Env.MetricsPort)
// 	srv := &http.Server{
// 		Handler:      mux,
// 		Addr:         fmt.Sprintf(":%s", queueConfig.Env.MetricsPort),
// 		WriteTimeout: 10 * time.Second,
// 		ReadTimeout:  10 * time.Second,
// 	}

// 	if err := srv.ListenAndServe(); err != nil {
// 		logger.Errorf("Error starting Queue Metrics Server: %v\n", err)
// 		return
// 	}
// }

func run() (func(), error) {
	app, mux, cleanup := buildServer()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Start metrics server
	// go startMetricsServer()

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
