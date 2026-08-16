# User Activity Tracking Service

A lightweight Go REST API service backed by PostgreSQL that ingests user activity events, exposes filtered query endpoints, and runs a background aggregation job every 4 hours. Includes a minimal React web client for viewing events.

---

## Project Structure

```
user-activity-tracking-service/
├── cmd/
│   ├── api/
│   │   └── main.go                 # Service entrypoint (config, migrations, pool)
│   └── seed/
│       └── main.go                 # Database seeder script
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
├── frontend/                       # React + TypeScript + Tailwind client
│   ├── src/                        # Component hierarchy & state management
│   ├── Dockerfile                  # Multi-stage Nginx container build
│   └── vite.config.ts              # Vite config with dev proxy to API
├── Makefile                    # Make targets (run, seed, seed-local, test, docker-up, etc.)
├── .env.example                # Sample environment variables
├── docker-compose.yml          # Multi-container setup (Postgres + API + Frontend)
├── go.mod                      # Go module definition
└── README.md
```

---

## Getting Started

### 1. Prerequisites
- [Go 1.22+](https://golang.org/dl/)
- [Node.js 20+](https://nodejs.org/)
- [Docker & Docker Compose](https://www.docker.com/)

### 2. Environment Configuration
Copy `.env.example` to `.env`:
```bash
cp .env.example .env
```

### 3. Running with Docker Compose
To build and run all services (PostgreSQL, Go REST API, and Nginx React Frontend):
```bash
docker compose up --build
# or using Make:
make docker-up
```
- **Frontend UI**: `http://localhost:3000`
- **REST API**: `http://localhost:8080`
- **PostgreSQL**: `localhost:5433` (host mapping) / `5432` (internal network)

#### Seed Data in Docker
To populate the running Docker containers with ~100 realistic events and compute historical 4-hour aggregations:
```bash
docker compose exec api /app/seed
# or using Make:
make seed
```

### 4. Running Locally for Development

#### A. Start PostgreSQL
```bash
docker compose up -d postgres
```

#### B. Run the API Service
```bash
go run ./cmd/api/main.go
# or using Make:
make run
```
The service will automatically run pending SQL migrations and establish a PostgreSQL connection pool.

#### C. Run the Frontend Client
```bash
cd frontend
npm install
npm run dev
# or from root using Make:
make frontend
```
Open `http://localhost:5173` in your browser. API calls to `/api` and `/health` will be automatically proxied to `http://localhost:8080`.

#### D. Seed Realistic Mock Data Locally (Optional)
```bash
go run ./cmd/seed/main.go
# or using Make:
make seed-local
```

### 5. Running Tests
```bash
go test -v ./...
# or using Make:
make test
```
