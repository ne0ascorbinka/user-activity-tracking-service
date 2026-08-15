# User Activity Tracking Service

A lightweight Go REST API service backed by PostgreSQL that ingests user activity events, exposes filtered query endpoints, and runs a background aggregation job every 4 hours. Includes a minimal React web client for viewing events.

---

## Project Structure

```
user-activity-tracking-service/
├── cmd/
│   └── api/
│       └── main.go                 # Service entrypoint (config, migrations, pool)
├── internal/
│   ├── config/
│   │   ├── config.go               # Configuration loader and PostgreSQL DSN generator
│   │   └── config_test.go          # Config unit tests
│   ├── database/
│   │   ├── database.go             # pgxpool connection & healthcheck logic
│   │   └── migrate.go              # golang-migrate runner
│   ├── models/
│   │   ├── event.go                # Ingested event models
│   │   └── stat.go                 # Aggregated stats models
│   ├── repository/                 # Data access layer
│   ├── service/                    # Business logic layer
│   ├── handler/                    # HTTP handlers
│   └── worker/                     # Background aggregation worker
├── migrations/
│   ├── 000001_create_events_table.up.sql
│   ├── 000001_create_events_table.down.sql
│   ├── 000002_create_user_activity_stats_table.up.sql
│   ├── 000002_create_user_activity_stats_table.down.sql
│   ├── migrations.go               # Embeds SQL files into binary
│   └── migrations_test.go          # Migration embed tests
├── docs/
│   └── SPEC.md                     # Project specification
├── .env.example                    # Sample environment variables
├── docker-compose.yml              # PostgreSQL container setup
├── go.mod                          # Go module definition
└── README.md
```

---

## Getting Started

### 1. Prerequisites
- [Go 1.22+](https://golang.org/dl/)
- [Docker & Docker Compose](https://www.docker.com/)

### 2. Environment Configuration
Copy `.env.example` to `.env`:
```bash
cp .env.example .env
```

### 3. Start PostgreSQL
```bash
docker compose up -d postgres
```

### 4. Run the API Service
```bash
go run ./cmd/api/main.go
```
The service will automatically run pending SQL migrations and establish a PostgreSQL connection pool.

### 5. Running Tests
```bash
go test -v ./...
```
