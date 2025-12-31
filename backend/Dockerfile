# Build stage
FROM golang:1.24-alpine AS builder

ARG SERVICE_NAME

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/service ./services/${SERVICE_NAME}

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/service .

# Run as non-root user
RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8080 50051

CMD ["./service"]

