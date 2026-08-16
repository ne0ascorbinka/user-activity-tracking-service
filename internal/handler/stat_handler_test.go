package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"user-activity-tracking-service/internal/models"
	"user-activity-tracking-service/internal/service"
)

// mockStatService implements service.StatService for handler testing.
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

func TestStatHandler_ListStats(t *testing.T) {
	start := time.Date(2026, 8, 16, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 8, 16, 4, 0, 0, 0, time.UTC)

	sampleStats := []models.UserActivityStat{
		{
			ID:          1,
			UserID:      42,
			EventCount:  15,
			PeriodStart: start,
			PeriodEnd:   end,
			CreatedAt:   end,
			UpdatedAt:   end,
		},
	}

	tests := []struct {
		name               string
		url                string
		mockService        *mockStatService
		expectedStatusCode int
		checkResponse      func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "200 OK - No filters (returns defaults)",
			url:  "/api/v1/stats",
			mockService: &mockStatService{
				getStatsFunc: func(ctx context.Context, filter models.ListStatsFilter) ([]models.UserActivityStat, error) {
					if filter.Limit != 50 || filter.Offset != 0 {
						t.Errorf("expected default limit 50, offset 0; got %d, %d", filter.Limit, filter.Offset)
					}
					return sampleStats, nil
				},
			},
			expectedStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var res []models.UserActivityStat
				if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
					t.Fatalf("failed to decode response JSON: %v", err)
				}
				if len(res) != 1 || res[0].UserID != 42 || res[0].EventCount != 15 {
					t.Errorf("unexpected response body: %+v", res)
				}
			},
		},
		{
			name: "200 OK - Valid user_id, pagination and date range",
			url:  "/api/v1/stats?user_id=42&from=2026-08-16T00:00:00Z&to=2026-08-16T04:00:00Z&limit=10&offset=5",
			mockService: &mockStatService{
				getStatsFunc: func(ctx context.Context, filter models.ListStatsFilter) ([]models.UserActivityStat, error) {
					if filter.UserID == nil || *filter.UserID != 42 {
						t.Errorf("expected user_id 42, got %v", filter.UserID)
					}
					if filter.From == nil || !filter.From.Equal(start) {
						t.Errorf("expected from %v, got %v", start, filter.From)
					}
					if filter.To == nil || !filter.To.Equal(end) {
						t.Errorf("expected to %v, got %v", end, filter.To)
					}
					if filter.Limit != 10 || filter.Offset != 5 {
						t.Errorf("expected limit 10, offset 5; got %d, %d", filter.Limit, filter.Offset)
					}
					return sampleStats, nil
				},
			},
			expectedStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var res []models.UserActivityStat
				if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
					t.Fatalf("failed to decode response JSON: %v", err)
				}
				if len(res) != 1 {
					t.Errorf("expected 1 result, got %d", len(res))
				}
			},
		},
		{
			name: "400 Bad Request - Invalid user_id non-numeric",
			url:  "/api/v1/stats?user_id=abc",
			mockService: &mockStatService{},
			expectedStatusCode: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var errRes ErrorResponse
				_ = json.Unmarshal(rec.Body.Bytes(), &errRes)
				if errRes.Error == "" {
					t.Error("expected error message in response")
				}
			},
		},
		{
			name: "400 Bad Request - Invalid user_id <= 0",
			url:  "/api/v1/stats?user_id=-1",
			mockService: &mockStatService{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "400 Bad Request - Invalid from timestamp",
			url:  "/api/v1/stats?from=not-a-date",
			mockService: &mockStatService{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "400 Bad Request - Invalid to timestamp",
			url:  "/api/v1/stats?to=2026-08-99",
			mockService: &mockStatService{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "400 Bad Request - Invalid limit <= 0",
			url:  "/api/v1/stats?limit=0",
			mockService: &mockStatService{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "400 Bad Request - Invalid offset < 0",
			url:  "/api/v1/stats?offset=-5",
			mockService: &mockStatService{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "400 Bad Request - Service ErrInvalidDateRange",
			url:  "/api/v1/stats?from=2026-08-16T10:00:00Z&to=2026-08-16T08:00:00Z",
			mockService: &mockStatService{
				getStatsFunc: func(ctx context.Context, filter models.ListStatsFilter) ([]models.UserActivityStat, error) {
					return nil, service.ErrInvalidDateRange
				},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "500 Internal Server Error - Service failure",
			url:  "/api/v1/stats",
			mockService: &mockStatService{
				getStatsFunc: func(ctx context.Context, filter models.ListStatsFilter) ([]models.UserActivityStat, error) {
					return nil, errors.New("database failure")
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewStatHandler(tt.mockService)
			mux := http.NewServeMux()
			h.RegisterRoutes(mux)

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatusCode {
				t.Fatalf("expected status code %d, got %d. Body: %s", tt.expectedStatusCode, rec.Code, rec.Body.String())
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}
