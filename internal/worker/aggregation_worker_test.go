package worker

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"user-activity-tracking-service/internal/models"
)

type mockStatService struct {
	runAggregationFunc func(ctx context.Context, now time.Time, interval time.Duration) (int64, error)
	getStatsFunc       func(ctx context.Context, filter models.ListStatsFilter) ([]models.UserActivityStat, error)
}

func (m *mockStatService) RunAggregation(ctx context.Context, now time.Time, interval time.Duration) (int64, error) {
	if m.runAggregationFunc != nil {
		return m.runAggregationFunc(ctx, now, interval)
	}
	return 0, nil
}

func (m *mockStatService) GetStats(ctx context.Context, filter models.ListStatsFilter) ([]models.UserActivityStat, error) {
	if m.getStatsFunc != nil {
		return m.getStatsFunc(ctx, filter)
	}
	return nil, nil
}

func TestAggregationWorker_StartAndStop(t *testing.T) {
	var callCount atomic.Int64
	mockSvc := &mockStatService{
		runAggregationFunc: func(ctx context.Context, now time.Time, interval time.Duration) (int64, error) {
			callCount.Add(1)
			return 2, nil
		},
	}

	nullLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	worker := NewAggregationWorker(mockSvc, 20*time.Millisecond, nullLogger)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		worker.Start(ctx)
		close(done)
	}()

	// Allow initial run + at least one tick
	time.Sleep(70 * time.Millisecond)
	cancel()

	select {
	case <-done:
		// Clean exit
	case <-time.After(500 * time.Millisecond):
		t.Fatal("worker did not shut down in time upon context cancellation")
	}

	calls := callCount.Load()
	if calls < 2 {
		t.Errorf("expected at least 2 aggregation runs, got %d", calls)
	}
}

func TestAggregationWorker_ServiceErrorHandling(t *testing.T) {
	mockSvc := &mockStatService{
		runAggregationFunc: func(ctx context.Context, now time.Time, interval time.Duration) (int64, error) {
			return 0, errors.New("database connection failed")
		},
	}

	nullLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	worker := NewAggregationWorker(mockSvc, 20*time.Millisecond, nullLogger)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		worker.Start(ctx)
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case <-done:
		// Clean exit despite service errors
	case <-time.After(500 * time.Millisecond):
		t.Fatal("worker did not shut down in time upon context cancellation")
	}
}

func TestAggregationWorker_PreCancelledContext(t *testing.T) {
	var callCount atomic.Int64
	mockSvc := &mockStatService{
		runAggregationFunc: func(ctx context.Context, now time.Time, interval time.Duration) (int64, error) {
			callCount.Add(1)
			return 0, nil
		},
	}

	nullLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	worker := NewAggregationWorker(mockSvc, 20*time.Millisecond, nullLogger)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan struct{})
	go func() {
		worker.Start(ctx)
		close(done)
	}()

	select {
	case <-done:
		// Exited immediately
	case <-time.After(500 * time.Millisecond):
		t.Fatal("worker did not exit immediately with cancelled context")
	}

	if callCount.Load() != 0 {
		t.Errorf("expected 0 calls for pre-cancelled context, got %d", callCount.Load())
	}
}
