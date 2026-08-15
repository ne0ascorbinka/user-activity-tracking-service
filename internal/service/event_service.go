package service

import (
	"context"
	"errors"
	"fmt"

	"user-activity-tracking-service/internal/models"
	"user-activity-tracking-service/internal/repository"
)

var (
	ErrInvalidDateRange = errors.New("'from' timestamp cannot be after 'to' timestamp")
)

// EventService defines business logic operations for events.
type EventService interface {
	IngestEvent(ctx context.Context, req models.IngestEventRequest) (*models.Event, error)
	GetEvents(ctx context.Context, filter models.ListEventsFilter) ([]models.Event, error)
}

// eventService implements EventService.
type eventService struct {
	repo repository.EventRepository
}

// NewEventService creates a new EventService instance.
func NewEventService(repo repository.EventRepository) EventService {
	return &eventService{
		repo: repo,
	}
}

// IngestEvent validates the incoming event request and persists it.
func (s *eventService) IngestEvent(ctx context.Context, req models.IngestEventRequest) (*models.Event, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	event := &models.Event{
		UserID:   req.UserID,
		Action:   req.Action,
		Metadata: req.Metadata,
	}

	created, err := s.repo.Create(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("service failed to create event: %w", err)
	}

	return created, nil
}

// GetEvents retrieves events filtered by criteria with pagination bounds.
func (s *eventService) GetEvents(ctx context.Context, filter models.ListEventsFilter) ([]models.Event, error) {
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

	events, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("service failed to list events: %w", err)
	}

	if events == nil {
		events = make([]models.Event, 0)
	}

	return events, nil
}
