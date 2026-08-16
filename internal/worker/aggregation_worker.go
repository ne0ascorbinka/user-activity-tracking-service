package worker

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"user-activity-tracking-service/internal/service"
)

// AggregationWorker periodically computes and persists aggregated user activity statistics.
type AggregationWorker struct {
	service  service.StatService
	interval time.Duration
	logger   *slog.Logger
}

// NewAggregationWorker creates a new AggregationWorker instance.
func NewAggregationWorker(svc service.StatService, interval time.Duration, logger *slog.Logger) *AggregationWorker {
	if logger == nil {
		logger = slog.Default()
	}
	return &AggregationWorker{
		service:  svc,
		interval: interval,
		logger:   logger,
	}
}

// Start begins periodic execution of the aggregation job until the context is cancelled.
func (w *AggregationWorker) Start(ctx context.Context) {
	w.logger.Info("starting aggregation worker", slog.Duration("interval", w.interval))

	// Initial run upon worker startup
	w.runOnce(ctx, time.Now())

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("aggregation worker shutting down cleanly")
			return
		case tickTime := <-ticker.C:
			w.runOnce(ctx, tickTime)
		}
	}
}

// runOnce executes a single aggregation job pass with the specified reference time.
func (w *AggregationWorker) runOnce(ctx context.Context, refTime time.Time) {
	if ctx.Err() != nil {
		return
	}

	w.logger.Info("running background aggregation job", slog.Time("trigger_time", refTime))
	rowsAffected, err := w.service.RunAggregation(ctx, refTime, w.interval)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			w.logger.Info("aggregation job cancelled during shutdown")
			return
		}
		w.logger.Error("aggregation job failed", slog.String("error", err.Error()))
		return
	}

	w.logger.Info("aggregation job completed successfully",
		slog.Int64("records_upserted", rowsAffected),
		slog.Duration("interval", w.interval),
	)
}
