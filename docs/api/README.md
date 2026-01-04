# API Documentation

Money Control API is a RESTful API built with Go and Gin framework. All endpoints are prefixed with `/api/v1`.

## Base URL

```
http://localhost:9080/api/v1
```

## Authentication

Most endpoints require authentication. Include the JWT token in the `Authorization` header:

```
Authorization: Bearer <your_access_token>
```

## Table of Contents

- [Authentication](./auth.md)
- [Accounts](./accounts.md)
- [Transactions](./transactions.md)
- [Assets](./assets.md)
- [Insights](./insights.md)
- [Currency](./currency.md)

## Health Endpoints

Health endpoints are available at the root level (not under `/api/v1`).

### Liveness Probe

```http
GET /health
```

Returns service running status. Use for Kubernetes `livenessProbe`.

**Response:**
```json
{
  "status": "UP",
  "service": "gateway",
  "version": "1.0.0",
  "timestamp": "2024-01-04T12:00:00Z"
}
```

### Readiness Probe

```http
GET /ready
```

Checks all dependencies (Redis, databases). Use for Kubernetes `readinessProbe`.

**Response (Healthy - 200):**
```json
{
  "status": "UP",
  "version": "1.0.0",
  "checks": {
    "redis": {
      "name": "redis",
      "status": "UP",
      "message": "Redis connection is healthy"
    }
  },
  "timestamp": "2024-01-04T12:00:00Z"
}
```

**Response (Unhealthy - 503):**
```json
{
  "status": "DOWN",
  "checks": {
    "redis": {
      "status": "DOWN",
      "message": "connection refused"
    }
  }
}
```

### Health Details

```http
GET /health/details
```

Full health report including circuit breaker status.

### Circuit Breaker Status

```http
GET /health/circuits
```

Returns status of all circuit breakers.

**Response:**
```json
{
  "accounts-service": {
    "name": "accounts-service",
    "state": "closed",
    "failures": 0,
    "successes": 10
  }
}
```

Circuit breaker states:
- `closed` - Normal operation
- `open` - Failing fast (service down)
- `half-open` - Testing if service recovered

## Response Format

### Success Response

```json
{
  "success": true,
  "data": { ... }
}
```

### Error Response

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description"
  }
}
```

## HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 429 | Too Many Requests (rate limited) |
| 500 | Internal Server Error |
| 503 | Service Unavailable (circuit open) |

## Rate Limiting

API requests are rate limited to 100 requests per minute per IP address (when Redis is enabled).

## Pagination

List endpoints support pagination with the following query parameters:

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | int | 1 | Page number |
| `page_size` | int | 20 | Items per page (max 100) |

Response includes pagination info:

```json
{
  "data": [...],
  "total": 150,
  "page": 1,
  "page_size": 20,
  "total_pages": 8
}
```

