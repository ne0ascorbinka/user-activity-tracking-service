package models

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidUserID  = errors.New("user_id must be a positive integer greater than 0")
	ErrEmptyAction    = errors.New("action is required and cannot be empty")
	ErrActionTooLong  = errors.New("action cannot exceed 64 characters")
	ErrInvalidMetadata = errors.New("metadata must be a valid JSON object or value")
)

// Event represents an ingested activity event.
type Event struct {
	ID        int64           `json:"id"`
	UserID    int64           `json:"user_id"`
	Action    string          `json:"action"`
	Metadata  json.RawMessage `json:"metadata"`
	CreatedAt time.Time       `json:"created_at"`
}

// IngestEventRequest represents the incoming JSON payload for ingesting an event.
type IngestEventRequest struct {
	UserID   int64           `json:"user_id"`
	Action   string          `json:"action"`
	Metadata json.RawMessage `json:"metadata"`
}

// Validate checks if the ingestion request is valid.
func (r *IngestEventRequest) Validate() error {
	if r.UserID <= 0 {
		return ErrInvalidUserID
	}

	trimmedAction := strings.TrimSpace(r.Action)
	if trimmedAction == "" {
		return ErrEmptyAction
	}
	if len(trimmedAction) > 64 {
		return ErrActionTooLong
	}
	r.Action = trimmedAction

	if len(r.Metadata) == 0 || string(r.Metadata) == "null" {
		r.Metadata = json.RawMessage("{}")
	} else if !json.Valid(r.Metadata) {
		return ErrInvalidMetadata
	}

	return nil
}

// ListEventsFilter defines query parameters for filtering events.
type ListEventsFilter struct {
	UserID *int64
	From   *time.Time
	To     *time.Time
	Limit  int
	Offset int
}
