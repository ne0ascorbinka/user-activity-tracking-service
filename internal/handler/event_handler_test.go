package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"user-activity-tracking-service/internal/models"
	"user-activity-tracking-service/internal/service"
)

// mockEventService implements service.EventService for testing.
type mockEventService struct {
	ingestEventFunc func(ctx context.Context, req models.IngestEventRequest) (*models.Event, error)
	getEventsFunc   func(ctx context.Context, filter models.ListEventsFilter) ([]models.Event, error)
}

func (m *mockEventService) IngestEvent(ctx context.Context, req models.IngestEventRequest) (*models.Event, error) {
	if m.ingestEventFunc != nil {
		return m.ingestEventFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockEventService) GetEvents(ctx context.Context, filter models.ListEventsFilter) ([]models.Event, error) {
	if m.getEventsFunc != nil {
		return m.getEventsFunc(ctx, filter)
	}
	return nil, nil
}

func TestEventHandler_IngestEvent(t *testing.T) {
	fixedTime := time.Date(2026, 8, 15, 18, 0, 0, 0, time.UTC)

	tests := []struct {
		name               string
		reqBody            string
		mockService        *mockEventService
		expectedStatusCode int
		checkResponse      func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:    "201 Created - Valid payload with metadata",
			reqBody: `{"user_id": 42, "action": "page_view", "metadata": {"page": "/home"}}`,
			mockService: &mockEventService{
				ingestEventFunc: func(ctx context.Context, req models.IngestEventRequest) (*models.Event, error) {
					if req.UserID != 42 || req.Action != "page_view" {
						return nil, errors.New("unexpected payload")
					}
					return &models.Event{
						ID:        1,
						UserID:    req.UserID,
						Action:    req.Action,
						Metadata:  req.Metadata,
						CreatedAt: fixedTime,
					}, nil
				},
			},
			expectedStatusCode: http.StatusCreated,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var ev models.Event
				if err := json.Unmarshal(rec.Body.Bytes(), &ev); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if ev.ID != 1 || ev.UserID != 42 || ev.Action != "page_view" {
					t.Errorf("unexpected event response: %+v", ev)
				}
			},
		},
		{
			name:    "201 Created - Valid payload without metadata",
			reqBody: `{"user_id": 100, "action": "login"}`,
			mockService: &mockEventService{
				ingestEventFunc: func(ctx context.Context, req models.IngestEventRequest) (*models.Event, error) {
					return &models.Event{
						ID:        2,
						UserID:    req.UserID,
						Action:    req.Action,
						Metadata:  json.RawMessage("{}"),
						CreatedAt: fixedTime,
					}, nil
				},
			},
			expectedStatusCode: http.StatusCreated,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var ev models.Event
				if err := json.Unmarshal(rec.Body.Bytes(), &ev); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if ev.UserID != 100 || ev.Action != "login" {
					t.Errorf("unexpected response: %+v", ev)
				}
			},
		},
		{
			name:               "400 Bad Request - Missing / zero user_id",
			reqBody:            `{"user_id": 0, "action": "page_view"}`,
			mockService:        &mockEventService{},
			expectedStatusCode: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var errResp ErrorResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}
				if !strings.Contains(errResp.Error, "user_id") {
					t.Errorf("expected user_id error, got: %s", errResp.Error)
				}
			},
		},
		{
			name:    "400 Bad Request - Service validation returns ErrInvalidUserID",
			reqBody: `{"user_id": -5, "action": "page_view"}`,
			mockService: &mockEventService{
				ingestEventFunc: func(ctx context.Context, req models.IngestEventRequest) (*models.Event, error) {
					return nil, models.ErrInvalidUserID
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var errResp ErrorResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}
				if errResp.Error != models.ErrInvalidUserID.Error() {
					t.Errorf("expected %v, got %v", models.ErrInvalidUserID, errResp.Error)
				}
			},
		},
		{
			name:    "400 Bad Request - Empty action",
			reqBody: `{"user_id": 42, "action": "  "}`,
			mockService: &mockEventService{
				ingestEventFunc: func(ctx context.Context, req models.IngestEventRequest) (*models.Event, error) {
					return nil, models.ErrEmptyAction
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var errResp ErrorResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}
				if errResp.Error != models.ErrEmptyAction.Error() {
					t.Errorf("expected %v, got %v", models.ErrEmptyAction, errResp.Error)
				}
			},
		},
		{
			name:    "400 Bad Request - Action exceeding 64 characters",
			reqBody: `{"user_id": 42, "action": "` + strings.Repeat("a", 65) + `"}`,
			mockService: &mockEventService{
				ingestEventFunc: func(ctx context.Context, req models.IngestEventRequest) (*models.Event, error) {
					return nil, models.ErrActionTooLong
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var errResp ErrorResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}
				if errResp.Error != models.ErrActionTooLong.Error() {
					t.Errorf("expected %v, got %v", models.ErrActionTooLong, errResp.Error)
				}
			},
		},
		{
			name:               "400 Bad Request - Malformed JSON",
			reqBody:            `{"user_id": 42, "action": `,
			mockService:        &mockEventService{},
			expectedStatusCode: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var errResp ErrorResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}
				if !strings.Contains(errResp.Error, "invalid JSON") {
					t.Errorf("expected invalid JSON error, got: %s", errResp.Error)
				}
			},
		},
		{
			name:               "400 Bad Request - Unknown field",
			reqBody:            `{"user_id": 42, "action": "click", "unknown": "value"}`,
			mockService:        &mockEventService{},
			expectedStatusCode: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var errResp ErrorResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}
				if !strings.Contains(errResp.Error, "unknown field") {
					t.Errorf("expected unknown field error, got: %s", errResp.Error)
				}
			},
		},
		{
			name:               "413 Payload Too Large - Body > 64 KB",
			reqBody:            `{"user_id": 42, "action": "test", "metadata": {"large": "` + strings.Repeat("x", 70*1024) + `"}}`,
			mockService:        &mockEventService{},
			expectedStatusCode: http.StatusRequestEntityTooLarge,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var errResp ErrorResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}
				if !strings.Contains(errResp.Error, "64 KB") {
					t.Errorf("expected 64 KB limit error, got: %s", errResp.Error)
				}
			},
		},
		{
			name:    "500 Internal Server Error - Service failure",
			reqBody: `{"user_id": 42, "action": "page_view"}`,
			mockService: &mockEventService{
				ingestEventFunc: func(ctx context.Context, req models.IngestEventRequest) (*models.Event, error) {
					return nil, errors.New("db down")
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var errResp ErrorResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}
				if errResp.Error != "failed to ingest event" {
					t.Errorf("expected internal error message, got: %s", errResp.Error)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewEventHandler(tt.mockService)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/events", bytes.NewBufferString(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			h.IngestEvent(rec, req)

			if rec.Code != tt.expectedStatusCode {
				t.Fatalf("expected status code %d, got %d. Body: %s", tt.expectedStatusCode, rec.Code, rec.Body.String())
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}

func TestEventHandler_ListEvents(t *testing.T) {
	fixedTime := time.Date(2026, 8, 15, 18, 0, 0, 0, time.UTC)

	tests := []struct {
		name               string
		url                string
		mockService        *mockEventService
		expectedStatusCode int
		checkResponse      func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "200 OK - No query parameters",
			url:  "/api/v1/events",
			mockService: &mockEventService{
				getEventsFunc: func(ctx context.Context, filter models.ListEventsFilter) ([]models.Event, error) {
					if filter.Limit != 50 || filter.Offset != 0 {
						t.Errorf("expected default pagination (50, 0), got (%d, %d)", filter.Limit, filter.Offset)
					}
					return []models.Event{
						{
							ID:        1,
							UserID:    42,
							Action:    "page_view",
							Metadata:  json.RawMessage(`{"page": "/home"}`),
							CreatedAt: fixedTime,
						},
					}, nil
				},
			},
			expectedStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var events []models.Event
				if err := json.Unmarshal(rec.Body.Bytes(), &events); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(events) != 1 || events[0].ID != 1 {
					t.Errorf("unexpected events returned: %+v", events)
				}
			},
		},
		{
			name: "200 OK - Valid filters with user_id, from, to, limit, offset",
			url:  "/api/v1/events?user_id=42&from=2026-08-01T00:00:00Z&to=2026-08-15T23:59:59Z&limit=25&offset=10",
			mockService: &mockEventService{
				getEventsFunc: func(ctx context.Context, filter models.ListEventsFilter) ([]models.Event, error) {
					if filter.UserID == nil || *filter.UserID != 42 {
						t.Errorf("expected UserID 42, got %v", filter.UserID)
					}
					if filter.From == nil || filter.To == nil {
						t.Errorf("expected from/to timestamps")
					}
					if filter.Limit != 25 || filter.Offset != 10 {
						t.Errorf("expected pagination (25, 10), got (%d, %d)", filter.Limit, filter.Offset)
					}
					return []models.Event{}, nil
				},
			},
			expectedStatusCode: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var events []models.Event
				if err := json.Unmarshal(rec.Body.Bytes(), &events); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(events) != 0 {
					t.Errorf("expected empty array, got %d items", len(events))
				}
			},
		},
		{
			name:               "400 Bad Request - Invalid user_id non-numeric",
			url:                "/api/v1/events?user_id=abc",
			mockService:        &mockEventService{},
			expectedStatusCode: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var errResp ErrorResponse
				_ = json.Unmarshal(rec.Body.Bytes(), &errResp)
				if !strings.Contains(errResp.Error, "user_id") {
					t.Errorf("expected user_id error, got: %s", errResp.Error)
				}
			},
		},
		{
			name:               "400 Bad Request - Invalid user_id <= 0",
			url:                "/api/v1/events?user_id=0",
			mockService:        &mockEventService{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "400 Bad Request - Invalid 'from' timestamp",
			url:                "/api/v1/events?from=invalid-date",
			mockService:        &mockEventService{},
			expectedStatusCode: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var errResp ErrorResponse
				_ = json.Unmarshal(rec.Body.Bytes(), &errResp)
				if !strings.Contains(errResp.Error, "RFC3339") {
					t.Errorf("expected RFC3339 error, got: %s", errResp.Error)
				}
			},
		},
		{
			name:               "400 Bad Request - Invalid 'to' timestamp",
			url:                "/api/v1/events?to=2026-99-99",
			mockService:        &mockEventService{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "400 Bad Request - Invalid limit negative",
			url:                "/api/v1/events?limit=-1",
			mockService:        &mockEventService{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "400 Bad Request - Invalid offset negative",
			url:                "/api/v1/events?offset=-5",
			mockService:        &mockEventService{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "400 Bad Request - 'from' after 'to' date range error",
			url:  "/api/v1/events?from=2026-08-15T00:00:00Z&to=2026-08-01T00:00:00Z",
			mockService: &mockEventService{
				getEventsFunc: func(ctx context.Context, filter models.ListEventsFilter) ([]models.Event, error) {
					return nil, service.ErrInvalidDateRange
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var errResp ErrorResponse
				_ = json.Unmarshal(rec.Body.Bytes(), &errResp)
				if errResp.Error != service.ErrInvalidDateRange.Error() {
					t.Errorf("expected date range error, got: %s", errResp.Error)
				}
			},
		},
		{
			name: "500 Internal Server Error - Service failure",
			url:  "/api/v1/events",
			mockService: &mockEventService{
				getEventsFunc: func(ctx context.Context, filter models.ListEventsFilter) ([]models.Event, error) {
					return nil, errors.New("db error")
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewEventHandler(tt.mockService)
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()

			h.ListEvents(rec, req)

			if rec.Code != tt.expectedStatusCode {
				t.Fatalf("expected status code %d, got %d. Body: %s", tt.expectedStatusCode, rec.Code, rec.Body.String())
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}
