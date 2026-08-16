package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"user-activity-tracking-service/internal/models"
	"user-activity-tracking-service/internal/repository"
)

var (
	ErrInvalidAggregationInterval = errors.New("aggregation interval must be greater than 0")
)

// StatService defines business logic operations for user activity statistics.
type StatService interface {
	RunAggregation(ctx context.Context, now time.Time, interval time.Duration) (int64, error)
	GetStats(ctx context.Context, filter models.ListStatsFilter) ([]models.UserActivityStat, error)
}

// statService implements StatService.
type statService struct {
	repo repository.StatRepository
}

// NewStatService creates a new StatService instance.
func NewStatService(repo repository.StatRepository) StatService {
	return &statService{
		repo: repo,
	}
}

// RunAggregation computes the aggregation window [now - interval, now] and executes upsert.
func (s *statService) RunAggregation(ctx context.Context, now time.Time, interval time.Duration) (int64, error) {
	if interval <= 0 {
		return 0, ErrInvalidAggregationInterval
	}

	periodStart := now.Add(-interval)
	periodEnd := now

	rowsAffected, err := s.repo.AggregateAndUpsert(ctx, periodStart, periodEnd)
	if err != nil {
		return 0, fmt.Errorf("service failed to run aggregation: %w", err)
	}

	return rowsAffected, nil
}

// GetStats retrieves user activity stats filtered by criteria with pagination bounds.
func (s *statService) GetStats(ctx context.Context, filter models.ListStatsFilter) ([]models.UserActivityStat, error) {
	if filter.UserID != nil && *filter.UserID <= 0 {
		return nil, models.ErrInvalidUserID
	}

	if filter.From != nil && filter.To != nil && filter.From.After(*filter.To) {
		return nil, ErrInvalidDateRange
	}

	if filter.Limit <= 0 {
		filter.Limit = 50
	} else if filter.Limit > 100 {
		filter.Limit = 100
	}

	if filter.Offset < 0 {
		filter.Offset = 0
	}

	stats, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("service failed to list stats: %w", err)
	}

	if stats == nil {
		stats = make([]models.UserActivityStat, 0)
	}

	return stats, nil
}
