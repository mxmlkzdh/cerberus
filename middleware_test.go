package cerberus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock implementation of RateLimiter
type MockRateLimiter struct {
	IsAllowedFunction func(*http.Request) (bool, error)
}

func (rl *MockRateLimiter) IsAllowed(r *http.Request) (bool, error) {
	return rl.IsAllowedFunction(r)
}

// Test middleware allows requests within the limit
func TestMiddlewareAllowsWithinLimit(t *testing.T) {
	mockLimiter := &MockRateLimiter{
		IsAllowedFunction: func(r *http.Request) (bool, error) {
			return true, nil
		},
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	middleware := Middleware(mockLimiter, handler)
	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", rr.Code)
	}
}

// Test middleware blocks requests exceeding the limit
func TestMiddlewareBlocksExceededLimit(t *testing.T) {
	mockLimiter := &MockRateLimiter{
		IsAllowedFunction: func(r *http.Request) (bool, error) {
			return false, nil
		},
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	middleware := Middleware(mockLimiter, handler)
	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("expected status Too Many Requests; got %v", rr.Code)
	}
}

// Test middleware returns an HTTP 500 on error
func TestMiddlewareHandlesError(t *testing.T) {
	mockLimiter := &MockRateLimiter{
		IsAllowedFunction: func(r *http.Request) (bool, error) {
			return false, fmt.Errorf("rate limiter error")
		},
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := Middleware(mockLimiter, handler)
	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status Internal Server Error; got %v", rr.Code)
	}
}
