# User Activity Tracking Service

A lightweight, robust Go REST API service backed by PostgreSQL that ingests high-throughput user activity events, provides filtered query endpoints, and executes scheduled background aggregations (every 4 hours). Includes a modern React web dashboard for inspecting events and aggregated metrics in real time.

---

## ⚡ Quickstart (Run Full Stack)

You can launch the entire stack (**PostgreSQL + Go Backend + React Frontend**) with a single command:

```bash
docker compose up --build
```

### Accessing the Services

| Service | URL / Port | Description |
| :--- | :--- | :--- |
| **Frontend Dashboard** | [http://localhost:3000](http://localhost:3000) | Interactive React UI for filtering events and viewing activity |
| **REST API** | [http://localhost:8080](http://localhost:8080) | Go HTTP server (`/api/v1/events`, `/health`) |
| **Grafana** | [http://localhost:3001](http://localhost:3001) | Observability UI & log explorer (admin / admin) |
| **Loki** | `http://localhost:3100` | Log aggregation engine |
| **PostgreSQL** | `localhost:5433` (Host) / `5432` (Internal) | Database (`activity_tracker`) |

> [!TIP]
> **Populate Demo Data:** Run the following command in another terminal while containers are running to seed realistic mock data and compute sample 4-hour aggregates:
> ```bash
> docker compose exec api /app/seed
> # or using Make:
> make seed
> ```

---

## 📊 Observability & Log Monitoring (Grafana + Loki + Promtail)

Container logs across the stack are automatically collected via Promtail, forwarded into Loki, and visualized in Grafana.

1. Open **Grafana** in your browser at [http://localhost:3001](http://localhost:3001).
2. Log in with the default credentials:
   - **Username**: `admin`
   - **Password**: `admin`
3. Navigate to **Explore** in the left sidebar (or press `g` then `e`).
4. Ensure the **Loki** data source is selected in the dropdown.
5. In the query editor, enter one of the following LogQL expressions to stream logs in real time:
   ```logql
   # View all logs from the Go REST API container:
   {container_name="activity_tracker_api"}

   # Or using the app label:
   {app="api"}

   # Filter specifically for errors or worker executions:
   {app="api"} |= "aggregation"
   ```

---

## 🏗 Architecture & Features

- **High-Performance Go Backend**: Built with standard library `net/http` routing, structured JSON logging, and connection pooling via `jackc/pgx/v5`.
- **Database Migrations**: Automatic SQL schema migrations applied on application startup using `golang-migrate`.
- **Background Worker**: In-process ticker worker computing aggregate activity stats over sliding 4-hour windows (configurable via `AGGREGATION_INTERVAL`).
- **React Frontend**: Clean, responsive dashboard built with Vite, TypeScript, and Tailwind CSS.
- **Dockerized**: Production-ready multi-stage Docker builds with health checks and volume persistence.

---

## 📡 API Reference & Examples

### 1. Ingest Activity Event

Record an action performed by a user.

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 42,
    "action": "page_view",
    "metadata": {
      "page": "/dashboard",
      "ip": "192.168.1.1"
    }
  }'
```

**Response (`201 Created`):**
```json
{
  "id": 1,
  "user_id": 42,
  "action": "page_view",
  "metadata": {
    "ip": "192.168.1.1",
    "page": "/dashboard"
  },
  "created_at": "2026-08-16T12:00:00Z"
}
```

---

### 2. Retrieve Filtered Events

Query stored events with optional filtering by `user_id`, `from`, and `to` timestamps (ISO 8601 / RFC3339).

```bash
# Retrieve all events
curl http://localhost:8080/api/v1/events

# Filter by user_id
curl "http://localhost:8080/api/v1/events?user_id=42"

# Filter by user_id and date range
curl "http://localhost:8080/api/v1/events?user_id=42&from=2026-08-01T00:00:00Z&to=2026-08-16T23:59:59Z"
```

---

### 3. Health Check

```bash
curl http://localhost:8080/health
```

**Response (`200 OK`):**
```json
{
  "status": "ok",
  "database": "connected"
}
```

---

## 💻 Local Development Workflow

If you prefer running components locally outside of Docker containers:

### Prerequisites
- [Go 1.22+](https://golang.org/dl/)
- [Node.js 20+](https://nodejs.org/)
- [Docker](https://www.docker.com/) (for PostgreSQL)

### 1. Configure Environment
```bash
cp .env.example .env
```

### 2. Start PostgreSQL Database
```bash
docker compose up -d postgres
```

### 3. Run Backend API
```bash
go run ./cmd/api/main.go
# or: make run
```
*Migrations will execute automatically on startup.*

### 4. Run Frontend Client
```bash
cd frontend
npm install
npm run dev
# or from project root: make frontend
```
*Frontend dev server starts at `http://localhost:5173` with automatic proxying to `http://localhost:8080`.*

### 5. Seed Local Database (Optional)
```bash
go run ./cmd/seed/main.go
# or: make seed-local
```

### 6. Run Tests
```bash
go test -v ./...
# or: make test
```

---

## 📂 Project Structure

```
user-activity-tracking-service/
├── cmd/
│   ├── api/main.go                 # Service entrypoint (config, migrations, pool)
│   └── seed/main.go                # Database seeder script
├── internal/
│   ├── config/                     # Configuration loader & PostgreSQL DSN builder
│   ├── database/                   # pgxpool connection & migration runner
│   ├── handler/                    # HTTP endpoints & request validators
│   ├── models/                     # Event & Stat domain models
│   ├── repository/                 # PostgreSQL query layer
│   ├── service/                    # Business logic layer
│   └── worker/                     # 4-hour periodic aggregation worker
├── migrations/                     # Embedded SQL migration files
├── docs/                           # Project specifications & documentation
├── frontend/                       # Vite + React + TypeScript web dashboard
├── docker-compose.yml              # PostgreSQL, API, and Frontend orchestration
├── Makefile                        # Shortcuts for development commands
└── README.md
```
