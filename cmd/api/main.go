package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"user-activity-tracking-service/internal/config"
	"user-activity-tracking-service/internal/database"
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

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	logger.Info("service scaffolding and database setup ready (press Ctrl+C to stop)")
	<-sigChan

	logger.Info("shutting down user-activity-tracking-service...")
}
