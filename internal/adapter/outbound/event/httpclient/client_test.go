package httpclient

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientOpensCircuitAfterFailures(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, err := New(Config{
		BaseURL: server.URL,
		CircuitBreaker: CircuitBreakerSettings{
			FailureThreshold: 2,
			OpenStateTimeout: time.Minute,
		},
	})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	request := Request{Method: http.MethodGet, Path: "/external"}

	if _, err := client.Do(context.Background(), request); err == nil {
		t.Fatal("expected first request to fail")
	}
	if _, err := client.Do(context.Background(), request); err == nil {
		t.Fatal("expected second request to fail")
	}
	if _, err := client.Do(context.Background(), request); !errors.Is(err, ErrCircuitBreakerOpen) {
		t.Fatalf("expected open circuit error, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 external calls before circuit opened, got %d", calls)
	}
}

func TestClientHalfOpenClosesAfterSuccessfulProbe(t *testing.T) {
	now := time.Now()
	shouldFail := true

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if shouldFail {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := New(Config{
		BaseURL: server.URL,
		CircuitBreaker: CircuitBreakerSettings{
			FailureThreshold: 1,
			SuccessThreshold: 1,
			OpenStateTimeout: time.Second,
		},
	})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	client.now = func() time.Time { return now }

	request := Request{Method: http.MethodGet, Path: "/external"}

	if _, err := client.Do(context.Background(), request); err == nil {
		t.Fatal("expected request to fail")
	}
	if _, err := client.Do(context.Background(), request); !errors.Is(err, ErrCircuitBreakerOpen) {
		t.Fatalf("expected open circuit error, got %v", err)
	}

	now = now.Add(2 * time.Second)
	shouldFail = false

	if _, err := client.Do(context.Background(), request); err != nil {
		t.Fatalf("expected half-open probe to succeed, got %v", err)
	}
	if _, err := client.Do(context.Background(), request); err != nil {
		t.Fatalf("expected closed circuit request to succeed, got %v", err)
	}
}
