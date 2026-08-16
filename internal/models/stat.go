package models

import (
	"time"
)

// UserActivityStat represents aggregated statistics for a user within a time interval.
type UserActivityStat struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	EventCount  int       `json:"event_count"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ListStatsFilter defines query parameters for filtering aggregate statistics.
type ListStatsFilter struct {
	UserID *int64
	From   *time.Time
	To     *time.Time
	Limit  int
	Offset int
}

