package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"user-activity-tracking-service/internal/models"
)

type mockStatRepository struct {
	aggregateAndUpsertFunc func(ctx context.Context, periodStart, periodEnd time.Time) (int64, error)
	listFunc               func(ctx context.Context, filter models.ListStatsFilter) ([]models.UserActivityStat, error)
}

func (m *mockStatRepository) AggregateAndUpsert(ctx context.Context, periodStart, periodEnd time.Time) (int64, error) {
	if m.aggregateAndUpsertFunc != nil {
		return m.aggregateAndUpsertFunc(ctx, periodStart, periodEnd)
	}
	return 0, nil
}

func (m *mockStatRepository) List(ctx context.Context, filter models.ListStatsFilter) ([]models.UserActivityStat, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, filter)
	}
	return nil, nil
}

func TestStatService_RunAggregation(t *testing.T) {
	now := time.Date(2026, 8, 16, 12, 0, 0, 0, time.UTC)
	interval := 4 * time.Hour

	tests := []struct {
		name         string
		now          time.Time
		interval     time.Duration
		mockRepo     *mockStatRepository
		expectedRows int64
		expectedErr  error
		checkWindow  func(t *testing.T, start, end time.Time)
	}{
		{
			name:     "Successful aggregation",
			now:      now,
			interval: interval,
			mockRepo: &mockStatRepository{
				aggregateAndUpsertFunc: func(ctx context.Context, periodStart, periodEnd time.Time) (int64, error) {
					expectedStart := now.Add(-interval)
					if !periodStart.Equal(expectedStart) || !periodEnd.Equal(now) {
						t.Errorf("expected window [%v, %v], got [%v, %v]", expectedStart, now, periodStart, periodEnd)
					}
					return 5, nil
				},
			},
			expectedRows: 5,
			expectedErr:  nil,
		},
		{
			name:        "Invalid interval <= 0",
			now:         now,
			interval:    0,
			mockRepo:    &mockStatRepository{},
			expectedErr: ErrInvalidAggregationInterval,
		},
		{
			name:     "Repository error propagation",
			now:      now,
			interval: interval,
			mockRepo: &mockStatRepository{
				aggregateAndUpsertFunc: func(ctx context.Context, periodStart, periodEnd time.Time) (int64, error) {
					return 0, errors.New("db query failed")
				},
			},
			expectedRows: 0,
			expectedErr:  errors.New("service failed to run aggregation: db query failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewStatService(tt.mockRepo)
			rows, err := svc.RunAggregation(context.Background(), tt.now, tt.interval)

			if tt.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.expectedErr)
				}
				if errors.Is(tt.expectedErr, ErrInvalidAggregationInterval) {
					if !errors.Is(err, tt.expectedErr) {
						t.Errorf("expected error %v, got %v", tt.expectedErr, err)
					}
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if rows != tt.expectedRows {
					t.Errorf("expected rows affected %d, got %d", tt.expectedRows, rows)
				}
			}
		})
	}
}

func TestStatService_GetStats(t *testing.T) {
	now := time.Now()
	past := now.Add(-4 * time.Hour)
	validUserID := int64(42)
	invalidUserID := int64(-5)

	tests := []struct {
		name        string
		filter      models.ListStatsFilter
		mockRepo    *mockStatRepository
		expectedErr error
		checkFilter func(t *testing.T, filter models.ListStatsFilter)
	}{
		{
			name: "Valid filter and default normalization",
			filter: models.ListStatsFilter{
				UserID: &validUserID,
				Limit:  0,
				Offset: -5,
			},
			mockRepo: &mockStatRepository{
				listFunc: func(ctx context.Context, filter models.ListStatsFilter) ([]models.UserActivityStat, error) {
					return []models.UserActivityStat{
						{
							ID:          1,
							UserID:      42,
							EventCount:  10,
							PeriodStart: past,
							PeriodEnd:   now,
							CreatedAt:   now,
							UpdatedAt:   now,
						},
					}, nil
				},
			},
			checkFilter: func(t *testing.T, filter models.ListStatsFilter) {
				if filter.Limit != 50 || filter.Offset != 0 {
					t.Errorf("expected normalized limit=50, offset=0; got limit=%d, offset=%d", filter.Limit, filter.Offset)
				}
			},
		},
		{
			name: "Limit upper bound clamping",
			filter: models.ListStatsFilter{
				Limit: 250,
			},
			mockRepo: &mockStatRepository{
				listFunc: func(ctx context.Context, filter models.ListStatsFilter) ([]models.UserActivityStat, error) {
					if filter.Limit != 100 {
						t.Errorf("expected limit clamped to 100, got %d", filter.Limit)
					}
					return nil, nil
				},
			},
		},
		{
			name: "Invalid user_id <= 0",
			filter: models.ListStatsFilter{
				UserID: &invalidUserID,
			},
			mockRepo:    &mockStatRepository{},
			expectedErr: models.ErrInvalidUserID,
		},
		{
			name: "Invalid date range: from after to",
			filter: models.ListStatsFilter{
				From: &now,
				To:   &past,
			},
			mockRepo:    &mockStatRepository{},
			expectedErr: ErrInvalidDateRange,
		},
		{
			name: "Repository error",
			filter: models.ListStatsFilter{
				UserID: &validUserID,
			},
			mockRepo: &mockStatRepository{
				listFunc: func(ctx context.Context, filter models.ListStatsFilter) ([]models.UserActivityStat, error) {
					return nil, errors.New("db error")
				},
			},
			expectedErr: errors.New("service failed to list stats: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewStatService(tt.mockRepo)
			stats, err := svc.GetStats(context.Background(), tt.filter)

			if tt.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.expectedErr)
				}
				if errors.Is(tt.expectedErr, models.ErrInvalidUserID) || errors.Is(tt.expectedErr, ErrInvalidDateRange) {
					if !errors.Is(err, tt.expectedErr) {
						t.Errorf("expected error %v, got %v", tt.expectedErr, err)
					}
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if stats == nil {
					t.Fatal("expected non-nil stats slice")
				}
			}
		})
	}
}
