# API Reference

Base URL: `http://localhost:9080/api/v1`

## Authentication

Include JWT token in Authorization header:

```
Authorization: Bearer <access_token>
```

## Endpoints

- [Authentication](./auth.md)
- [Accounts](./accounts.md)
- [Transactions](./transactions.md)
- [Assets](./assets.md)
- [Insights](./insights.md)
- [Currency](./currency.md)

## Health Endpoints

Health endpoints are at root level (not `/api/v1`).

| Endpoint | Method | Description | Response |
|----------|--------|-------------|----------|
| `/health` | GET | Liveness probe | `{"status": "UP"}` |
| `/ready` | GET | Readiness probe | `{"status": "UP", "checks": {...}}` |
| `/health/circuits` | GET | Circuit breaker status | `{"service": {"state": "closed"}}` |

### Readiness Response Codes

| Code | Meaning |
|------|---------|
| 200 | All dependencies healthy |
| 503 | One or more dependencies unhealthy |

## Response Format

### Success

```json
{
  "success": true,
  "data": {}
}
```

### Error

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Description"
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
| 429 | Rate Limited |
| 500 | Internal Error |
| 503 | Service Unavailable |

## Rate Limiting

100 requests per minute per IP. Headers returned:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1704369600
```

## Pagination

| Parameter | Type | Default | Max |
|-----------|------|---------|-----|
| `page` | int | 1 | - |
| `page_size` | int | 20 | 100 |

Response:

```json
{
  "data": [],
  "total": 150,
  "page": 1,
  "page_size": 20
}
```
