package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"user-activity-tracking-service/internal/config"
	"user-activity-tracking-service/internal/database"
	"user-activity-tracking-service/internal/repository"
	"user-activity-tracking-service/internal/service"
)

type SeedEvent struct {
	UserID    int64
	Action    string
	Metadata  map[string]any
	CreatedAt time.Time
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("starting database seeder...")

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// 1. Run database migrations if not already applied
	if err := database.RunMigrations(cfg.DSN(), logger); err != nil {
		logger.Error("failed to run database migrations", slog.String("error", err.Error()))
		os.Exit(1)
	}

	ctx := context.Background()

	// 2. Connect to database
	pool, err := database.NewPool(ctx, cfg, logger)
	if err != nil {
		logger.Error("failed to initialize database connection pool", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	// 3. Generate realistic mock events
	events := generateSeedEvents()
	logger.Info("generated mock events", slog.Int("count", len(events)))

	// 4. Batch insert events into PostgreSQL
	batch := &pgx.Batch{}
	for _, ev := range events {
		metaJSON, err := json.Marshal(ev.Metadata)
		if err != nil {
			logger.Error("failed to marshal metadata", slog.String("error", err.Error()))
			os.Exit(1)
		}

		batch.Queue(
			"INSERT INTO events (user_id, action, metadata, created_at) VALUES ($1, $2, $3, $4)",
			ev.UserID,
			ev.Action,
			metaJSON,
			ev.CreatedAt,
		)
	}

	br := pool.SendBatch(ctx, batch)
	for i := 0; i < len(events); i++ {
		if _, err := br.Exec(); err != nil {
			logger.Error("failed executing batch insert", slog.Int("index", i), slog.String("error", err.Error()))
			_ = br.Close()
			os.Exit(1)
		}
	}
	if err := br.Close(); err != nil {
		logger.Error("failed to close batch results", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("successfully inserted mock events", slog.Int("total_events", len(events)))

	// 5. Pre-aggregate historical 4-hour statistics across the 48-hour range
	statRepo := repository.NewPostgresStatRepository(pool)
	statService := service.NewStatService(statRepo)

	now := time.Now().UTC()
	windowDuration := 4 * time.Hour
	totalAggregatedRows := int64(0)

	// Step backward 12 windows (48 hours)
	for i := 12; i >= 0; i-- {
		windowEnd := now.Add(-time.Duration(i) * windowDuration)
		rows, err := statService.RunAggregation(ctx, windowEnd, windowDuration)
		if err != nil {
			logger.Warn("historical aggregation warning",
				slog.Time("period_end", windowEnd),
				slog.String("error", err.Error()),
			)
			continue
		}
		totalAggregatedRows += rows
	}

	logger.Info("historical periodic aggregations completed",
		slog.Int64("total_stat_records_upserted", totalAggregatedRows),
	)

	fmt.Printf("\n✓ Successfully seeded %d events across 5 users and computed historical 4-hour stats!\n", len(events))
}

func generateSeedEvents() []SeedEvent {
	userIDs := []int64{42, 52, 101, 7, 13}
	browsers := []string{"Chrome", "Safari", "Firefox", "Edge"}
	pages := []string{"/home", "/pricing", "/features", "/dashboard", "/settings", "/docs/api", "/checkout"}
	buttons := []string{"cta-start-free", "btn-pricing-annual", "btn-dark-mode-toggle", "btn-download-sdk", "btn-checkout-confirm"}
	items := []struct {
		id    string
		price float64
	}{
		{"starter_tier_monthly", 29.00},
		{"pro_tier_annual", 240.00},
		{"enterprise_custom", 999.00},
		{"api_credits_pack_100k", 49.50},
	}

	rng := rand.New(rand.NewSource(42)) // Deterministic seed for reproducible data

	var events []SeedEvent
	now := time.Now().UTC()

	// Distribute across 48 hours (~100 events)
	for hoursAgo := 47; hoursAgo >= 0; hoursAgo-- {
		// 1 to 4 events per hour
		eventsThisHour := rng.Intn(3) + 1

		for e := 0; e < eventsThisHour; e++ {
			userID := userIDs[rng.Intn(len(userIDs))]
			minuteOffset := rng.Intn(60)
			secondOffset := rng.Intn(60)
			eventTime := now.Add(-time.Duration(hoursAgo)*time.Hour + time.Duration(minuteOffset)*time.Minute + time.Duration(secondOffset)*time.Second)

			// Action types & tailored metadata
			actionRoll := rng.Float64()
			var action string
			var metadata map[string]any

			switch {
			case actionRoll < 0.40:
				action = "PAGE_VIEW"
				metadata = map[string]any{
					"page":     pages[rng.Intn(len(pages))],
					"browser":  browsers[rng.Intn(len(browsers))],
					"referrer": fmt.Sprintf("https://referrer-%d.org", rng.Intn(5)+1),
				}
			case actionRoll < 0.65:
				action = "BUTTON_CLICK"
				metadata = map[string]any{
					"button_id": buttons[rng.Intn(len(buttons))],
					"screen_x":  rng.Intn(1920),
					"screen_y":  rng.Intn(1080),
				}
			case actionRoll < 0.80:
				action = "LOGIN"
				metadata = map[string]any{
					"method":     []string{"password", "oauth_github", "oauth_google"}[rng.Intn(3)],
					"ip_address": fmt.Sprintf("192.168.%d.%d", rng.Intn(254)+1, rng.Intn(254)+1),
					"mfa":        rng.Intn(2) == 1,
				}
			case actionRoll < 0.92:
				action = "CHECKOUT"
				item := items[rng.Intn(len(items))]
				metadata = map[string]any{
					"item_id":  item.id,
					"price":    item.price,
					"currency": "USD",
					"status":   "completed",
				}
			default:
				action = "LOGOUT"
				metadata = map[string]any{
					"session_duration_sec": rng.Intn(3600) + 120,
					"reason":               []string{"user_action", "timeout"}[rng.Intn(2)],
				}
			}

			events = append(events, SeedEvent{
				UserID:    userID,
				Action:    action,
				Metadata:  metadata,
				CreatedAt: eventTime,
			})
		}
	}

	return events
}
