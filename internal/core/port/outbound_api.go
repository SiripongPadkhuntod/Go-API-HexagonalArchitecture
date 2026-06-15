package port

import (
	"context"
	"errors"
)

var ErrCircuitBreakerOpen = errors.New("circuit breaker is open")

type OutboundAPIClient interface {
	Do(ctx context.Context, request OutboundAPIRequest) (OutboundAPIResponse, error)
}

type OutboundAPIRequest struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    []byte
}

type OutboundAPIResponse struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}
