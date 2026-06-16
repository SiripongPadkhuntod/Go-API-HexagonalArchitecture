package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"hexagonalarchitecture/internal/core/domain"
	"hexagonalarchitecture/internal/core/port"
)

var ErrCircuitBreakerOpen = errors.New("circuit breaker is open")

type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    []byte
}

type Response struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

type userCreatedEvent struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var _ port.UserEventPublisher = (*Client)(nil)

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

func (c *Client) PublishUserCreated(ctx context.Context, user domain.User) error {
	payload, err := json.Marshal(userCreatedEvent{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		return err
	}

	_, err = c.Do(ctx, Request{
		Method: http.MethodPost,
		Path:   "/users/events",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: payload,
	})
	return err
}

func (c *Client) Do(ctx context.Context, request Request) (Response, error) {
	now := c.now()
	if err := c.breaker.beforeRequest(now); err != nil {
		return Response{}, err
	}

	response, err := c.do(ctx, request)
	c.breaker.afterRequest(err == nil, c.now())
	return response, err
}

func (c *Client) do(ctx context.Context, request Request) (Response, error) {
	method := strings.TrimSpace(request.Method)
	if method == "" {
		method = http.MethodGet
	}

	endpoint := c.baseURL.ResolveReference(&url.URL{Path: request.Path})
	httpRequest, err := http.NewRequestWithContext(ctx, method, endpoint.String(), bytes.NewReader(request.Body))
	if err != nil {
		return Response{}, fmt.Errorf("create outbound request: %w", err)
	}

	for key, value := range request.Headers {
		httpRequest.Header.Set(key, value)
	}

	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return Response{}, fmt.Errorf("send outbound request: %w", err)
	}
	defer httpResponse.Body.Close()

	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return Response{}, fmt.Errorf("read outbound response: %w", err)
	}

	response := Response{
		StatusCode: httpResponse.StatusCode,
		Headers:    httpResponse.Header,
		Body:       body,
	}
	if httpResponse.StatusCode >= http.StatusInternalServerError {
		return response, fmt.Errorf("outbound API returned status %d", httpResponse.StatusCode)
	}

	return response, nil
}
