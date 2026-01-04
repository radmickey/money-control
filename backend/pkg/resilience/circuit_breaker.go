package resilience

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Circuit breaker states
type State int

const (
	StateClosed State = iota // Normal operation, requests pass through
	StateOpen                // Circuit is open, requests fail fast
	StateHalfOpen            // Testing if service recovered
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Errors
var (
	ErrCircuitOpen    = errors.New("circuit breaker is open")
	ErrTooManyFails   = errors.New("too many failures")
	ErrRequestTimeout = errors.New("request timeout")
)

// CircuitBreakerConfig holds circuit breaker configuration
type CircuitBreakerConfig struct {
	Name             string        // Name for logging/metrics
	MaxFailures      int           // Max failures before opening circuit
	Timeout          time.Duration // How long circuit stays open
	HalfOpenRequests int           // Requests allowed in half-open state
	SuccessThreshold int           // Successes needed to close circuit
}

// DefaultConfig returns default circuit breaker configuration
func DefaultConfig(name string) CircuitBreakerConfig {
	return CircuitBreakerConfig{
		Name:             name,
		MaxFailures:      5,
		Timeout:          30 * time.Second,
		HalfOpenRequests: 3,
		SuccessThreshold: 2,
	}
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config CircuitBreakerConfig

	mu               sync.RWMutex
	state            State
	failures         int
	successes        int
	lastFailure      time.Time
	halfOpenRequests int
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	if config.MaxFailures <= 0 {
		config.MaxFailures = 5
	}
	if config.Timeout <= 0 {
		config.Timeout = 30 * time.Second
	}
	if config.HalfOpenRequests <= 0 {
		config.HalfOpenRequests = 3
	}
	if config.SuccessThreshold <= 0 {
		config.SuccessThreshold = 2
	}

	return &CircuitBreaker{
		config: config,
		state:  StateClosed,
	}
}

// Execute runs the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func(context.Context) error) error {
	if !cb.AllowRequest() {
		return ErrCircuitOpen
	}

	// Execute the function
	err := fn(ctx)

	// Record result
	if err != nil {
		cb.RecordFailure()
		return err
	}

	cb.RecordSuccess()
	return nil
}

// AllowRequest checks if a request should be allowed
func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return true

	case StateOpen:
		// Check if timeout has passed
		if time.Since(cb.lastFailure) > cb.config.Timeout {
			cb.toHalfOpen()
			return true
		}
		return false

	case StateHalfOpen:
		// Allow limited requests in half-open state
		if cb.halfOpenRequests < cb.config.HalfOpenRequests {
			cb.halfOpenRequests++
			return true
		}
		return false

	default:
		return false
	}
}

// RecordSuccess records a successful request
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		cb.failures = 0 // Reset failures on success

	case StateHalfOpen:
		cb.successes++
		if cb.successes >= cb.config.SuccessThreshold {
			cb.toClosed()
		}
	}
}

// RecordFailure records a failed request
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.lastFailure = time.Now()

	switch cb.state {
	case StateClosed:
		cb.failures++
		if cb.failures >= cb.config.MaxFailures {
			cb.toOpen()
		}

	case StateHalfOpen:
		cb.toOpen()
	}
}

// State transitions
func (cb *CircuitBreaker) toOpen() {
	cb.state = StateOpen
	cb.successes = 0
	cb.halfOpenRequests = 0
}

func (cb *CircuitBreaker) toHalfOpen() {
	cb.state = StateHalfOpen
	cb.successes = 0
	cb.halfOpenRequests = 0
}

func (cb *CircuitBreaker) toClosed() {
	cb.state = StateClosed
	cb.failures = 0
	cb.successes = 0
	cb.halfOpenRequests = 0
}

// GetState returns current circuit breaker state
func (cb *CircuitBreaker) GetState() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// GetName returns circuit breaker name
func (cb *CircuitBreaker) GetName() string {
	return cb.config.Name
}

// Stats returns circuit breaker statistics
func (cb *CircuitBreaker) Stats() map[string]interface{} {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return map[string]interface{}{
		"name":       cb.config.Name,
		"state":      cb.state.String(),
		"failures":   cb.failures,
		"successes":  cb.successes,
		"lastFailure": cb.lastFailure,
	}
}

// CircuitBreakerManager manages multiple circuit breakers
type CircuitBreakerManager struct {
	mu       sync.RWMutex
	breakers map[string]*CircuitBreaker
}

// NewCircuitBreakerManager creates a new circuit breaker manager
func NewCircuitBreakerManager() *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// Get returns a circuit breaker by name, creating it if necessary
func (m *CircuitBreakerManager) Get(name string) *CircuitBreaker {
	m.mu.RLock()
	cb, exists := m.breakers[name]
	m.mu.RUnlock()

	if exists {
		return cb
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check after acquiring write lock
	if cb, exists = m.breakers[name]; exists {
		return cb
	}

	cb = NewCircuitBreaker(DefaultConfig(name))
	m.breakers[name] = cb
	return cb
}

// GetWithConfig returns a circuit breaker with custom config
func (m *CircuitBreakerManager) GetWithConfig(config CircuitBreakerConfig) *CircuitBreaker {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cb, exists := m.breakers[config.Name]; exists {
		return cb
	}

	cb := NewCircuitBreaker(config)
	m.breakers[config.Name] = cb
	return cb
}

// AllStats returns stats for all circuit breakers
func (m *CircuitBreakerManager) AllStats() map[string]map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make(map[string]map[string]interface{})
	for name, cb := range m.breakers {
		stats[name] = cb.Stats()
	}
	return stats
}

// GlobalManager is the default circuit breaker manager
var GlobalManager = NewCircuitBreakerManager()

// Execute is a convenience function using the global manager
func Execute(ctx context.Context, name string, fn func(context.Context) error) error {
	return GlobalManager.Get(name).Execute(ctx, fn)
}

