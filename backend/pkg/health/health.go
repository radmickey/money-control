package health

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Status represents health status
type Status string

const (
	StatusUp      Status = "UP"
	StatusDown    Status = "DOWN"
	StatusUnknown Status = "UNKNOWN"
)

// Check represents a single health check
type Check struct {
	Name      string                 `json:"name"`
	Status    Status                 `json:"status"`
	Message   string                 `json:"message,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// HealthReport represents overall system health
type HealthReport struct {
	Status    Status           `json:"status"`
	Timestamp time.Time        `json:"timestamp"`
	Checks    map[string]Check `json:"checks"`
	Version   string           `json:"version,omitempty"`
}

// Checker is a function that performs a health check
type Checker func(ctx context.Context) Check

// HealthChecker manages health checks
type HealthChecker struct {
	mu       sync.RWMutex
	checkers map[string]Checker
	version  string
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(version string) *HealthChecker {
	return &HealthChecker{
		checkers: make(map[string]Checker),
		version:  version,
	}
}

// Register adds a health checker
func (h *HealthChecker) Register(name string, checker Checker) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checkers[name] = checker
}

// Check performs all health checks
func (h *HealthChecker) Check(ctx context.Context) HealthReport {
	h.mu.RLock()
	defer h.mu.RUnlock()

	report := HealthReport{
		Status:    StatusUp,
		Timestamp: time.Now(),
		Checks:    make(map[string]Check),
		Version:   h.version,
	}

	for name, checker := range h.checkers {
		check := checker(ctx)
		check.Name = name
		check.Timestamp = time.Now()
		report.Checks[name] = check

		if check.Status == StatusDown {
			report.Status = StatusDown
		}
	}

	return report
}

// Liveness returns liveness status (is the service running?)
func (h *HealthChecker) Liveness(ctx context.Context) Check {
	return Check{
		Name:      "liveness",
		Status:    StatusUp,
		Timestamp: time.Now(),
		Message:   "Service is running",
	}
}

// Readiness returns readiness status (can the service handle requests?)
func (h *HealthChecker) Readiness(ctx context.Context) HealthReport {
	return h.Check(ctx)
}

// DatabaseCheck creates a database health checker
func DatabaseCheck(pingFn func(ctx context.Context) error) Checker {
	return func(ctx context.Context) Check {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := pingFn(ctx); err != nil {
			return Check{
				Status:  StatusDown,
				Message: err.Error(),
			}
		}

		return Check{
			Status:  StatusUp,
			Message: "Database connection is healthy",
		}
	}
}

// RedisCheck creates a Redis health checker
func RedisCheck(pingFn func(ctx context.Context) error) Checker {
	return func(ctx context.Context) Check {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		if err := pingFn(ctx); err != nil {
			return Check{
				Status:  StatusDown,
				Message: err.Error(),
			}
		}

		return Check{
			Status:  StatusUp,
			Message: "Redis connection is healthy",
		}
	}
}

// GRPCServiceCheck creates a gRPC service health checker
func GRPCServiceCheck(serviceName string, conn *grpc.ClientConn) Checker {
	return func(ctx context.Context) Check {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		client := grpc_health_v1.NewHealthClient(conn)
		resp, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{
			Service: serviceName,
		})

		if err != nil {
			return Check{
				Status:  StatusDown,
				Message: err.Error(),
				Details: map[string]interface{}{
					"service": serviceName,
				},
			}
		}

		status := StatusDown
		if resp.Status == grpc_health_v1.HealthCheckResponse_SERVING {
			status = StatusUp
		}

		return Check{
			Status:  status,
			Message: resp.Status.String(),
			Details: map[string]interface{}{
				"service": serviceName,
			},
		}
	}
}

// GRPCHealthServer implements gRPC health checking protocol
type GRPCHealthServer struct {
	grpc_health_v1.UnimplementedHealthServer
	mu       sync.RWMutex
	services map[string]grpc_health_v1.HealthCheckResponse_ServingStatus
}

// NewGRPCHealthServer creates a new gRPC health server
func NewGRPCHealthServer() *GRPCHealthServer {
	return &GRPCHealthServer{
		services: make(map[string]grpc_health_v1.HealthCheckResponse_ServingStatus),
	}
}

// SetServingStatus sets the serving status of a service
func (s *GRPCHealthServer) SetServingStatus(service string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.services[service] = status
}

// SetServing marks a service as serving
func (s *GRPCHealthServer) SetServing(service string) {
	s.SetServingStatus(service, grpc_health_v1.HealthCheckResponse_SERVING)
}

// SetNotServing marks a service as not serving
func (s *GRPCHealthServer) SetNotServing(service string) {
	s.SetServingStatus(service, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
}

// Check implements the gRPC Health Check method
func (s *GRPCHealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if status, ok := s.services[req.Service]; ok {
		return &grpc_health_v1.HealthCheckResponse{
			Status: status,
		}, nil
	}

	// If service not found, check empty string (overall health)
	if req.Service == "" {
		// Default to serving if no services registered
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_SERVING,
		}, nil
	}

	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN,
	}, nil
}

// Watch implements the gRPC Health Watch method (streaming)
func (s *GRPCHealthServer) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	// Simple implementation - send current status and keep connection open
	s.mu.RLock()
	status := s.services[req.Service]
	s.mu.RUnlock()

	return stream.Send(&grpc_health_v1.HealthCheckResponse{
		Status: status,
	})
}

// Register registers the health server with a gRPC server
func (s *GRPCHealthServer) Register(srv *grpc.Server) {
	grpc_health_v1.RegisterHealthServer(srv, s)
}

