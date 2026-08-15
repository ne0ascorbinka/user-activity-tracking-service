package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"user-activity-tracking-service/internal/config"
)

// NewPool creates, configures, and verifies a PostgreSQL connection pool.
func NewPool(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database configuration: %w", err)
	}

	poolConfig.MaxConns = cfg.DBMaxConns
	poolConfig.MinConns = cfg.DBMinConns
	poolConfig.MaxConnLifetime = cfg.DBMaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.DBMaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection pool: %w", err)
	}

	// Retry pinging database with backoff
	const maxRetries = 10
	var pingErr error
	for i := 1; i <= maxRetries; i++ {
		pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		pingErr = pool.Ping(pingCtx)
		cancel()

		if pingErr == nil {
			logger.Info("connected to PostgreSQL successfully",
				slog.String("host", cfg.DBHost),
				slog.Int("port", cfg.DBPort),
				slog.String("database", cfg.DBName),
			)
			return pool, nil
		}

		logger.Warn("waiting for database connection...",
			slog.Int("attempt", i),
			slog.Int("max_attempts", maxRetries),
			slog.String("error", pingErr.Error()),
		)
		time.Sleep(1 * time.Second)
	}

	pool.Close()
	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, pingErr)
}
