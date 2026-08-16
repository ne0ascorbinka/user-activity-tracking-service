FROM golang:1.25-alpine AS builder

ENV GOTOOLCHAIN=auto

WORKDIR /app

# Install ca-certificates and git
RUN apk add --no-cache ca-certificates tzdata

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# Build statically compiled binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/api ./cmd/api

# Final runtime stage
FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/bin/api /app/api

EXPOSE 8080

ENTRYPOINT ["/app/api"]
