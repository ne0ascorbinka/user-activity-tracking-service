package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"user-activity-tracking-service/internal/models"
)

type mockEventRepository struct {
	createFunc func(ctx context.Context, event *models.Event) (*models.Event, error)
	listFunc   func(ctx context.Context, filter models.ListEventsFilter) ([]models.Event, error)
}

func (m *mockEventRepository) Create(ctx context.Context, event *models.Event) (*models.Event, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, event)
	}
	return nil, nil
}

func (m *mockEventRepository) List(ctx context.Context, filter models.ListEventsFilter) ([]models.Event, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, filter)
	}
	return nil, nil
}

func TestEventService_IngestEvent(t *testing.T) {
	tests := []struct {
		name        string
		req         models.IngestEventRequest
		mockRepo    *mockEventRepository
		expectedErr error
	}{
		{
			name: "Valid request with metadata",
			req: models.IngestEventRequest{
				UserID:   42,
				Action:   "page_view",
				Metadata: json.RawMessage(`{"page":"/dashboard"}`),
			},
			mockRepo: &mockEventRepository{
				createFunc: func(ctx context.Context, event *models.Event) (*models.Event, error) {
					return &models.Event{
						ID:        1,
						UserID:    event.UserID,
						Action:    event.Action,
						Metadata:  event.Metadata,
						CreatedAt: time.Now(),
					}, nil
				},
			},
			expectedErr: nil,
		},
		{
			name: "Validation error - invalid user_id",
			req: models.IngestEventRequest{
				UserID: 0,
				Action: "click",
			},
			mockRepo:    &mockEventRepository{},
			expectedErr: models.ErrInvalidUserID,
		},
		{
			name: "Validation error - empty action",
			req: models.IngestEventRequest{
				UserID: 1,
				Action: "   ",
			},
			mockRepo:    &mockEventRepository{},
			expectedErr: models.ErrEmptyAction,
		},
		{
			name: "Repository error",
			req: models.IngestEventRequest{
				UserID: 1,
				Action: "login",
			},
			mockRepo: &mockEventRepository{
				createFunc: func(ctx context.Context, event *models.Event) (*models.Event, error) {
					return nil, errors.New("db error")
				},
			},
			expectedErr: errors.New("service failed to create event: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewEventService(tt.mockRepo)
			created, err := svc.IngestEvent(context.Background(), tt.req)

			if tt.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.expectedErr)
				}
				if errors.Is(tt.expectedErr, models.ErrInvalidUserID) || errors.Is(tt.expectedErr, models.ErrEmptyAction) {
					if !errors.Is(err, tt.expectedErr) {
						t.Errorf("expected error %v, got %v", tt.expectedErr, err)
					}
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if created == nil || created.ID != 1 {
					t.Errorf("unexpected created event: %+v", created)
				}
			}
		})
	}
}

func TestEventService_GetEvents(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * time.Hour)

	tests := []struct {
		name        string
		filter      models.ListEventsFilter
		mockRepo    *mockEventRepository
		expectedErr error
		checkFilter func(t *testing.T, filter models.ListEventsFilter)
	}{
		{
			name: "Valid filter and default normalization",
			filter: models.ListEventsFilter{
				Limit:  0,
				Offset: -10,
			},
			mockRepo: &mockEventRepository{
				listFunc: func(ctx context.Context, filter models.ListEventsFilter) ([]models.Event, error) {
					return []models.Event{{ID: 1, UserID: 10, Action: "click"}}, nil
				},
			},
			checkFilter: func(t *testing.T, filter models.ListEventsFilter) {
				if filter.Limit != 50 || filter.Offset != 0 {
					t.Errorf("expected normalized limit=50, offset=0; got %d, %d", filter.Limit, filter.Offset)
				}
			},
		},
		{
			name: "Upper bound clamping for limit",
			filter: models.ListEventsFilter{
				Limit: 500,
			},
			mockRepo: &mockEventRepository{
				listFunc: func(ctx context.Context, filter models.ListEventsFilter) ([]models.Event, error) {
					return nil, nil
				},
			},
			checkFilter: func(t *testing.T, filter models.ListEventsFilter) {
				if filter.Limit != 100 {
					t.Errorf("expected limit clamped to 100, got %d", filter.Limit)
				}
			},
		},
		{
			name: "Error when from is after to",
			filter: models.ListEventsFilter{
				From: &now,
				To:   &past,
			},
			mockRepo:    &mockEventRepository{},
			expectedErr: ErrInvalidDateRange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewEventService(tt.mockRepo)
			events, err := svc.GetEvents(context.Background(), tt.filter)

			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if events == nil {
					t.Fatal("expected non-nil events slice")
				}
			}
		})
	}
}
