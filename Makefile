.PHONY: run seed seed-local test build frontend docker-up docker-down

# Run the API service locally
run:
	go run ./cmd/api/main.go

# Seed the database inside the running Docker container
seed:
	docker compose exec api /app/seed

# Seed the database from host
seed-local:
	go run ./cmd/seed/main.go

# Run Go tests
test:
	go test -v ./...

# Build API and Seeder binaries
build:
	go build -o ./bin/api ./cmd/api
	go build -o ./bin/seed ./cmd/seed

# Start the Vite React development server
frontend:
	cd frontend && npm run dev

# Build and run all services with Docker Compose
docker-up:
	docker compose up --build

# Stop Docker Compose containers
docker-down:
	docker compose down
