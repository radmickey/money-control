# Resilience Patterns

## Overview

| Pattern | Purpose |
|---------|---------|
| gRPC Retry | Automatic retry with exponential backoff |
| Circuit Breaker | Fail fast when downstream service is unhealthy |
| Health Checks | Kubernetes-ready liveness and readiness probes |
| Timeouts | Prevent hanging requests |

---

## gRPC Retry Policy

Configuration applied to all inter-service calls:

```json
{
  "timeout": "10s",
  "retryPolicy": {
    "maxAttempts": 3,
    "initialBackoff": "0.1s",
    "maxBackoff": "1s",
    "backoffMultiplier": 2.0,
    "retryableStatusCodes": ["UNAVAILABLE", "DEADLINE_EXCEEDED", "RESOURCE_EXHAUSTED"]
  }
}
```

**Retry timing:**

| Attempt | Backoff |
|---------|---------|
| 1 | 0ms |
| 2 | 100ms |
| 3 | 200ms |

**Non-retryable errors:** `NOT_FOUND`, `INVALID_ARGUMENT`, `PERMISSION_DENIED`, `UNAUTHENTICATED`

---

## Circuit Breaker

### State Diagram

```
CLOSED ──(5 failures)──> OPEN ──(30s)──> HALF-OPEN ──(2 successes)──> CLOSED
                           │                  │
                           └──────(failure)───┘
```

### Configuration

| Parameter | Default |
|-----------|---------|
| MaxFailures | 5 |
| Timeout | 30s |
| HalfOpenRequests | 3 |
| SuccessThreshold | 2 |

### Usage

```go
resp, err := resilience.Call(ctx, resilience.DefaultCallOptions("accounts-service"),
    func(ctx context.Context) (*accountspb.Account, error) {
        return client.CreateAccount(ctx, req)
    })
```

### Error Handling

```go
if errors.Is(err, resilience.ErrCircuitOpen) {
    return c.JSON(503, gin.H{"error": "Service temporarily unavailable", "retry_after": 30})
}
```

---

## Health Checks

### Endpoints

| Endpoint | HTTP Code | Purpose |
|----------|-----------|---------|
| `GET /health` | 200 | Liveness probe |
| `GET /ready` | 200/503 | Readiness probe |
| `GET /health/circuits` | 200 | Circuit breaker status |

### Response: `/health`

```json
{"status": "UP", "service": "gateway", "version": "1.0.0"}
```

### Response: `/ready`

```json
{
  "status": "UP",
  "checks": {
    "redis": {"status": "UP"}
  }
}
```

Returns `503 Service Unavailable` if any check fails.

### Response: `/health/circuits`

```json
{
  "accounts-service": {"state": "closed", "failures": 0},
  "currency-service": {"state": "open", "failures": 5}
}
```

---

## Timeouts

| Component | Value |
|-----------|-------|
| HTTP Read | 15s |
| HTTP Write | 15s |
| HTTP Idle | 60s |
| gRPC Call | 10s |

Custom timeout:

```go
opts := resilience.CallOptions{
    Timeout:     5 * time.Second,
    ServiceName: "accounts-service",
    UseBreaker:  true,
}
resp, err := resilience.Call(ctx, opts, fn)
```

---

## Keepalive

| Parameter | Value |
|-----------|-------|
| Ping Interval | 10s |
| Ping Timeout | 3s |
| PermitWithoutStream | true |

---

## Kubernetes Configuration

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 9080
  initialDelaySeconds: 5
  periodSeconds: 10
  failureThreshold: 3

readinessProbe:
  httpGet:
    path: /ready
    port: 9080
  initialDelaySeconds: 5
  periodSeconds: 5
  failureThreshold: 3
```
