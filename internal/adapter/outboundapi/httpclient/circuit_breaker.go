package httpclient

import (
	"sync"
	"time"

	"hexagonalarchitecture/internal/core/port"
)

type CircuitBreakerSettings struct {
	FailureThreshold uint32
	SuccessThreshold uint32
	OpenStateTimeout time.Duration
}

type circuitBreakerState string

const (
	circuitClosed   circuitBreakerState = "closed"
	circuitOpen     circuitBreakerState = "open"
	circuitHalfOpen circuitBreakerState = "half-open"
)

type circuitBreaker struct {
	mu sync.Mutex

	settings CircuitBreakerSettings
	state    circuitBreakerState
	openedAt time.Time

	consecutiveFailures  uint32
	consecutiveSuccesses uint32
}

func newCircuitBreaker(settings CircuitBreakerSettings) *circuitBreaker {
	if settings.FailureThreshold == 0 {
		settings.FailureThreshold = 3
	}
	if settings.SuccessThreshold == 0 {
		settings.SuccessThreshold = 1
	}
	if settings.OpenStateTimeout == 0 {
		settings.OpenStateTimeout = 30 * time.Second
	}

	return &circuitBreaker{
		settings: settings,
		state:    circuitClosed,
	}
}

func (b *circuitBreaker) beforeRequest(now time.Time) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state != circuitOpen {
		return nil
	}

	if now.Sub(b.openedAt) < b.settings.OpenStateTimeout {
		return port.ErrCircuitBreakerOpen
	}

	b.state = circuitHalfOpen
	b.consecutiveFailures = 0
	b.consecutiveSuccesses = 0
	return nil
}

func (b *circuitBreaker) afterRequest(success bool, now time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if success {
		b.onSuccess()
		return
	}

	b.onFailure(now)
}

func (b *circuitBreaker) onSuccess() {
	b.consecutiveFailures = 0

	if b.state != circuitHalfOpen {
		return
	}

	b.consecutiveSuccesses++
	if b.consecutiveSuccesses >= b.settings.SuccessThreshold {
		b.state = circuitClosed
		b.consecutiveSuccesses = 0
	}
}

func (b *circuitBreaker) onFailure(now time.Time) {
	b.consecutiveSuccesses = 0

	if b.state == circuitHalfOpen {
		b.open(now)
		return
	}

	b.consecutiveFailures++
	if b.consecutiveFailures >= b.settings.FailureThreshold {
		b.open(now)
	}
}

func (b *circuitBreaker) open(now time.Time) {
	b.state = circuitOpen
	b.openedAt = now
	b.consecutiveFailures = 0
	b.consecutiveSuccesses = 0
}
