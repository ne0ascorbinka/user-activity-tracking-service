package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"user-activity-tracking-service/internal/models"
)

// StatRepository defines data access methods for aggregated user activity stats.
type StatRepository interface {
	AggregateAndUpsert(ctx context.Context, periodStart, periodEnd time.Time) (int64, error)
	List(ctx context.Context, filter models.ListStatsFilter) ([]models.UserActivityStat, error)
}

// PostgresStatRepository implements StatRepository backed by PostgreSQL.
type PostgresStatRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresStatRepository creates a new PostgresStatRepository instance.
func NewPostgresStatRepository(pool *pgxpool.Pool) *PostgresStatRepository {
	return &PostgresStatRepository{pool: pool}
}

// AggregateAndUpsert calculates event counts per user within [periodStart, periodEnd] and upserts into user_activity_stats.
func (r *PostgresStatRepository) AggregateAndUpsert(ctx context.Context, periodStart, periodEnd time.Time) (int64, error) {
	query := `
		INSERT INTO user_activity_stats (user_id, event_count, period_start, period_end, created_at, updated_at)
		SELECT user_id, COUNT(*) AS event_count, $1 AS period_start, $2 AS period_end, NOW(), NOW()
		FROM events
		WHERE created_at >= $1 AND created_at <= $2
		GROUP BY user_id
		ON CONFLICT (user_id, period_start, period_end)
		DO UPDATE SET
			event_count = EXCLUDED.event_count,
			updated_at = NOW()
	`

	cmdTag, err := r.pool.Exec(ctx, query, periodStart, periodEnd)
	if err != nil {
		return 0, fmt.Errorf("failed to aggregate and upsert activity stats: %w", err)
	}

	return cmdTag.RowsAffected(), nil
}

// List queries user activity stats matching the given filter parameters.
func (r *PostgresStatRepository) List(ctx context.Context, filter models.ListStatsFilter) ([]models.UserActivityStat, error) {
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
		conditions = append(conditions, fmt.Sprintf("period_start >= $%d", argIndex))
		args = append(args, *filter.From)
		argIndex++
	}

	if filter.To != nil {
		conditions = append(conditions, fmt.Sprintf("period_end <= $%d", argIndex))
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
		SELECT id, user_id, event_count, period_start, period_end, created_at, updated_at
		FROM user_activity_stats
		%s
		ORDER BY period_start DESC, id DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query activity stats: %w", err)
	}
	defer rows.Close()

	stats := make([]models.UserActivityStat, 0)
	for rows.Next() {
		var stat models.UserActivityStat
		if err := rows.Scan(
			&stat.ID,
			&stat.UserID,
			&stat.EventCount,
			&stat.PeriodStart,
			&stat.PeriodEnd,
			&stat.CreatedAt,
			&stat.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan activity stat row: %w", err)
		}
		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during activity stats rows iteration: %w", err)
	}

	return stats, nil
}
