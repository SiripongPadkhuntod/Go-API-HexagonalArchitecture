package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"hexagonalarchitecture/internal/core/port"
)

type Config struct {
	BaseURL        string
	Timeout        time.Duration
	CircuitBreaker CircuitBreakerSettings
}

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	breaker    *circuitBreaker
	now        func() time.Time
}

func New(config Config) (*Client, error) {
	if strings.TrimSpace(config.BaseURL) == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	baseURL, err := url.Parse(config.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base URL: %w", err)
	}

	timeout := config.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		breaker: newCircuitBreaker(config.CircuitBreaker),
		now:     time.Now,
	}, nil
}

func (c *Client) Do(ctx context.Context, request port.OutboundAPIRequest) (port.OutboundAPIResponse, error) {
	now := c.now()
	if err := c.breaker.beforeRequest(now); err != nil {
		return port.OutboundAPIResponse{}, err
	}

	response, err := c.do(ctx, request)
	c.breaker.afterRequest(err == nil, c.now())
	return response, err
}

func (c *Client) do(ctx context.Context, request port.OutboundAPIRequest) (port.OutboundAPIResponse, error) {
	method := strings.TrimSpace(request.Method)
	if method == "" {
		method = http.MethodGet
	}

	endpoint := c.baseURL.ResolveReference(&url.URL{Path: request.Path})
	httpRequest, err := http.NewRequestWithContext(ctx, method, endpoint.String(), bytes.NewReader(request.Body))
	if err != nil {
		return port.OutboundAPIResponse{}, fmt.Errorf("create outbound request: %w", err)
	}

	for key, value := range request.Headers {
		httpRequest.Header.Set(key, value)
	}

	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return port.OutboundAPIResponse{}, fmt.Errorf("send outbound request: %w", err)
	}
	defer httpResponse.Body.Close()

	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return port.OutboundAPIResponse{}, fmt.Errorf("read outbound response: %w", err)
	}

	response := port.OutboundAPIResponse{
		StatusCode: httpResponse.StatusCode,
		Headers:    httpResponse.Header,
		Body:       body,
	}
	if httpResponse.StatusCode >= http.StatusInternalServerError {
		return response, fmt.Errorf("outbound API returned status %d", httpResponse.StatusCode)
	}

	return response, nil
}
