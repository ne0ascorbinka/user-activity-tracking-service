package models

import (
	"encoding/json"
	"time"
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
