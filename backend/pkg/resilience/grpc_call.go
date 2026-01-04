package resilience

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DefaultCallTimeout is the default timeout for gRPC calls
const DefaultCallTimeout = 10 * time.Second

// CallOptions holds options for a resilient gRPC call
type CallOptions struct {
	Timeout      time.Duration
	ServiceName  string
	UseBreaker   bool
}

// DefaultCallOptions returns default call options
func DefaultCallOptions(serviceName string) CallOptions {
	return CallOptions{
		Timeout:     DefaultCallTimeout,
		ServiceName: serviceName,
		UseBreaker:  true,
	}
}

// CallWithTimeout executes a gRPC call with timeout
func CallWithTimeout[T any](ctx context.Context, timeout time.Duration, fn func(context.Context) (T, error)) (T, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return fn(ctx)
}

// CallWithBreaker executes a gRPC call with circuit breaker protection
func CallWithBreaker[T any](ctx context.Context, serviceName string, fn func(context.Context) (T, error)) (T, error) {
	var result T
	var callErr error

	cb := GlobalManager.Get(serviceName)

	err := cb.Execute(ctx, func(ctx context.Context) error {
		result, callErr = fn(ctx)
		if callErr != nil {
			// Only count certain errors as failures
			if isRetryableError(callErr) {
				return callErr
			}
		}
		return nil
	})

	if err == ErrCircuitOpen {
		return result, status.Error(codes.Unavailable, "service temporarily unavailable (circuit open)")
	}

	return result, callErr
}

// Call executes a gRPC call with both timeout and circuit breaker
func Call[T any](ctx context.Context, opts CallOptions, fn func(context.Context) (T, error)) (T, error) {
	// Apply timeout
	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	// Apply circuit breaker if enabled
	if opts.UseBreaker {
		return CallWithBreaker(ctx, opts.ServiceName, fn)
	}

	return fn(ctx)
}

// isRetryableError checks if an error should trigger circuit breaker
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	st, ok := status.FromError(err)
	if !ok {
		return true // Unknown errors are considered retryable
	}

	switch st.Code() {
	case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted, codes.Aborted:
		return true
	case codes.NotFound, codes.InvalidArgument, codes.PermissionDenied, codes.Unauthenticated:
		return false // These are business logic errors, not infrastructure failures
	default:
		return false
	}
}

// MustCall is like Call but panics on circuit breaker open
// Use only in non-critical paths
func MustCall[T any](ctx context.Context, serviceName string, fn func(context.Context) (T, error)) (T, error) {
	return Call(ctx, DefaultCallOptions(serviceName), fn)
}

