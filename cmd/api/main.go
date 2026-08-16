package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"user-activity-tracking-service/internal/config"
	"user-activity-tracking-service/internal/database"
	"user-activity-tracking-service/internal/handler"
	"user-activity-tracking-service/internal/repository"
	"user-activity-tracking-service/internal/service"
	"user-activity-tracking-service/internal/worker"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("starting user-activity-tracking-service...")

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("configuration loaded",
		slog.Int("server_port", cfg.ServerPort),
		slog.String("db_host", cfg.DBHost),
		slog.Int("db_port", cfg.DBPort),
		slog.String("db_name", cfg.DBName),
		slog.Duration("aggregation_interval", cfg.AggregationInterval),
	)

	// Run database migrations
	if err := database.RunMigrations(cfg.DSN(), logger); err != nil {
		logger.Error("failed to run database migrations", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Create root context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database connection pool
	pool, err := database.NewPool(ctx, cfg, logger)
	if err != nil {
		logger.Error("failed to initialize database connection pool", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	logger.Info("database initialization and healthcheck completed successfully")

	// Initialize layers
	eventRepo := repository.NewPostgresEventRepository(pool)
	eventService := service.NewEventService(eventRepo)
	eventHandler := handler.NewEventHandler(eventService)

	statRepo := repository.NewPostgresStatRepository(pool)
	statService := service.NewStatService(statRepo)
	statHandler := handler.NewStatHandler(statService)

	// Router setup
	mux := http.NewServeMux()
	eventHandler.RegisterRoutes(mux)
	statHandler.RegisterRoutes(mux)

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Initialize and start background aggregation worker
	workerCtx, workerCancel := context.WithCancel(ctx)
	defer workerCancel()

	aggregationWorker := worker.NewAggregationWorker(statService, cfg.AggregationInterval, logger)
	go aggregationWorker.Start(workerCtx)

	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in background
	go func() {
		logger.Info("HTTP server listening", slog.String("addr", serverAddr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	logger.Info("shutting down HTTP server and background worker...")

	// Cancel background worker context
	workerCancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", slog.String("error", err.Error()))
	} else {
		logger.Info("HTTP server stopped gracefully")
	}

	logger.Info("service shutdown complete")
}
