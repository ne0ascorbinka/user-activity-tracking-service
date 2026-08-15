package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"user-activity-tracking-service/internal/models"
)

// EventRepository defines data access methods for events.
type EventRepository interface {
	Create(ctx context.Context, event *models.Event) (*models.Event, error)
	List(ctx context.Context, filter models.ListEventsFilter) ([]models.Event, error)
}

// PostgresEventRepository implements EventRepository backed by PostgreSQL.
type PostgresEventRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresEventRepository creates a new PostgresEventRepository instance.
func NewPostgresEventRepository(pool *pgxpool.Pool) *PostgresEventRepository {
	return &PostgresEventRepository{pool: pool}
}

// Create inserts a new event record into the database and returns the populated event.
func (r *PostgresEventRepository) Create(ctx context.Context, event *models.Event) (*models.Event, error) {
	query := `
		INSERT INTO events (user_id, action, metadata, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, user_id, action, metadata, created_at
	`

	var created models.Event
	err := r.pool.QueryRow(ctx, query, event.UserID, event.Action, event.Metadata).Scan(
		&created.ID,
		&created.UserID,
		&created.Action,
		&created.Metadata,
		&created.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert event: %w", err)
	}

	return &created, nil
}

// List queries events matching the given filter parameters.
func (r *PostgresEventRepository) List(ctx context.Context, filter models.ListEventsFilter) ([]models.Event, error) {
	var (
		conditions []string
		args       []any
		argIndex   = 1
	)

	if filter.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, *filter.UserID)
		argIndex++
	}

	if filter.From != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *filter.From)
		argIndex++
	}

	if filter.To != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *filter.To)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	} else if limit > 100 {
		limit = 100
	}

	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, action, metadata, created_at
		FROM events
		%s
		ORDER BY created_at DESC, id DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	events := make([]models.Event, 0)
	for rows.Next() {
		var ev models.Event
		if err := rows.Scan(&ev.ID, &ev.UserID, &ev.Action, &ev.Metadata, &ev.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan event row: %w", err)
		}
		events = append(events, ev)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return events, nil
}
