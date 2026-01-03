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
| 500 | Internal Server Error |

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

