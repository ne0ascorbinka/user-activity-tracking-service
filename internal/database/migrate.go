package database

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"user-activity-tracking-service/migrations"
)

// RunMigrations applies all pending database migrations.
func RunMigrations(dsn string, logger *slog.Logger) error {
	d, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to create iofs migration driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate instance: %w", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			logger.Warn("error closing migration source", slog.String("error", srcErr.Error()))
		}
		if dbErr != nil {
			logger.Warn("error closing migration db driver", slog.String("error", dbErr.Error()))
		}
	}()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("database schema is up to date (no migrations applied)")
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	logger.Info("database migrations applied successfully")
	return nil
}
