# Resilience Patterns

Money Control implements several resilience patterns to ensure high availability and graceful degradation under load.

## Overview

| Pattern | Purpose | Implementation |
|---------|---------|----------------|
| **gRPC Retry** | Handle transient failures | Automatic retry with exponential backoff |
| **Circuit Breaker** | Prevent cascading failures | Fail fast when service is unhealthy |
| **Health Checks** | Monitor service health | Kubernetes-ready probes |
| **Timeouts** | Prevent hanging requests | Configurable per-request timeouts |

## gRPC Retry Policy

All inter-service gRPC calls are configured with automatic retry:

```go
// Configuration
{
    "timeout": "10s",
    "retryPolicy": {
        "maxAttempts": 3,
        "initialBackoff": "0.1s",
        "maxBackoff": "1s",
        "backoffMultiplier": 2.0,
        "retryableStatusCodes": [
            "UNAVAILABLE",
            "DEADLINE_EXCEEDED",
            "RESOURCE_EXHAUSTED"
        ]
    }
}
```

### Retry Behavior

| Attempt | Backoff | Total Time |
|---------|---------|------------|
| 1 | 0ms | 0ms |
| 2 | 100ms | 100ms |
| 3 | 200ms | 300ms |

### Non-Retryable Errors

These errors are considered business logic errors and won't trigger retry:
- `NOT_FOUND` - Resource doesn't exist
- `INVALID_ARGUMENT` - Bad request data
- `PERMISSION_DENIED` - Authorization failed
- `UNAUTHENTICATED` - Authentication required

## Circuit Breaker

The circuit breaker prevents cascading failures when a downstream service is unhealthy.

### States

```
┌─────────┐     5 failures      ┌────────┐
│ CLOSED  │ ─────────────────→  │  OPEN  │
│ (Normal)│                     │(Reject)│
└────┬────┘                     └───┬────┘
     │                              │
     │ ◄──── 2 successes ─────┐     │ 30s timeout
     │                        │     │
     │                   ┌────┴─────▼────┐
     └───────────────────│  HALF-OPEN   │
                         │   (Testing)   │
                         └───────────────┘
```

### Configuration

| Parameter | Default | Description |
|-----------|---------|-------------|
| `MaxFailures` | 5 | Failures before opening circuit |
| `Timeout` | 30s | How long circuit stays open |
| `HalfOpenRequests` | 3 | Requests allowed in half-open state |
| `SuccessThreshold` | 2 | Successes needed to close circuit |

### Usage in Code

```go
import "github.com/radmickey/money-control/backend/pkg/resilience"

// Using the Call wrapper
resp, err := resilience.Call(ctx, resilience.DefaultCallOptions("accounts-service"),
    func(ctx context.Context) (*accountspb.Account, error) {
        return client.CreateAccount(ctx, req)
    })

// Handling circuit open
if errors.Is(err, resilience.ErrCircuitOpen) {
    // Service is temporarily unavailable
    return nil, status.Error(codes.Unavailable, "service temporarily unavailable")
}
```

### Monitoring Circuit Breakers

```bash
# Get all circuit breaker states
curl http://localhost:9080/health/circuits

# Example response
{
  "accounts-service": {
    "name": "accounts-service",
    "state": "closed",
    "failures": 0,
    "successes": 5
  },
  "currency-service": {
    "name": "currency-service",
    "state": "open",
    "failures": 5,
    "lastFailure": "2024-01-04T12:00:00Z"
  }
}
```

## Health Checks

### Endpoints

| Endpoint | Purpose | Kubernetes Probe |
|----------|---------|------------------|
| `GET /health` | Liveness check | `livenessProbe` |
| `GET /ready` | Readiness check | `readinessProbe` |
| `GET /health/details` | Full health report | Debugging |
| `GET /health/circuits` | Circuit breaker status | Monitoring |

### Liveness Probe

Checks if the service is running:

```bash
curl http://localhost:9080/health

# Response
{
  "status": "UP",
  "service": "gateway",
  "version": "1.0.0",
  "timestamp": "2024-01-04T12:00:00Z"
}
```

### Readiness Probe

Checks if all dependencies are healthy:

```bash
curl http://localhost:9080/ready

# Response (healthy)
{
  "status": "UP",
  "version": "1.0.0",
  "checks": {
    "redis": {
      "name": "redis",
      "status": "UP",
      "message": "Redis connection is healthy"
    }
  }
}

# Response (unhealthy) - returns HTTP 503
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

### Kubernetes Configuration

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gateway
spec:
  template:
    spec:
      containers:
      - name: gateway
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

## Timeouts

### Request Timeouts

| Component | Timeout | Description |
|-----------|---------|-------------|
| HTTP Read | 15s | Max time to read request |
| HTTP Write | 15s | Max time to write response |
| HTTP Idle | 60s | Keep-alive connection timeout |
| gRPC Call | 10s | Per-call timeout (configurable) |

### Custom Timeout

```go
// Use custom timeout
opts := resilience.CallOptions{
    Timeout:     5 * time.Second,  // Custom timeout
    ServiceName: "accounts-service",
    UseBreaker:  true,
}

resp, err := resilience.Call(ctx, opts, func(ctx context.Context) (*Response, error) {
    return client.FastOperation(ctx, req)
})
```

## Keepalive

gRPC connections are maintained with keepalive pings:

| Parameter | Value | Description |
|-----------|-------|-------------|
| Time | 10s | Ping interval when idle |
| Timeout | 3s | Wait for ping response |
| PermitWithoutStream | true | Ping even without active RPCs |

## Best Practices

### 1. Always Use Timeouts

```go
// ❌ Bad - no timeout
resp, err := client.Call(ctx, req)

// ✅ Good - with timeout
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()
resp, err := client.Call(ctx, req)
```

### 2. Handle Circuit Open

```go
// ✅ Return meaningful error to client
if errors.Is(err, resilience.ErrCircuitOpen) {
    return c.JSON(503, gin.H{
        "error": "Service temporarily unavailable",
        "retry_after": 30,
    })
}
```

### 3. Log Circuit State Changes

Monitor your logs for circuit breaker state changes to identify unhealthy services early.

## Metrics & Monitoring

For production, consider adding:
- Prometheus metrics for circuit breaker states
- Grafana dashboards for health visualization
- Alerting on circuit opens

```go
// Example: Prometheus metrics (future enhancement)
var circuitOpenTotal = prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Name: "circuit_breaker_open_total",
        Help: "Total number of circuit breaker opens",
    },
    []string{"service"},
)
```

