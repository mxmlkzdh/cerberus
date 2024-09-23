package cerberus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Mock implementation of AdvancedRateLimiter
type MockAdvancedRateLimiter struct {
	IsAllowedFunc        func(*http.Request) (bool, error)
	GetRateLimitDataFunc func(*http.Request) RateLimitData
}

func (m *MockAdvancedRateLimiter) IsAllowed(r *http.Request) (bool, error) {
	return m.IsAllowedFunc(r)
}

func (m *MockAdvancedRateLimiter) GetRateLimitData(r *http.Request) RateLimitData {
	return m.GetRateLimitDataFunc(r)
}

// Test allowing requests within the limit
func TestAdvancedMiddlewareAllowsWithinLimit(t *testing.T) {
	mockLimiter := &MockAdvancedRateLimiter{
		IsAllowedFunc: func(r *http.Request) (bool, error) {
			return true, nil
		},
		GetRateLimitDataFunc: func(r *http.Request) RateLimitData {
			return RateLimitData{
				Limit:      100,
				Remaining:  99,
				RetryAfter: 0,
			}
		},
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	middleware := AdvancedMiddleware(mockLimiter, handler)
	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", rr.Code)
	}
}

// Test blocking requests exceeding the limit
func TestAdvancedMiddlewareBlocksExceededLimit(t *testing.T) {
	mockLimiter := &MockAdvancedRateLimiter{
		IsAllowedFunc: func(r *http.Request) (bool, error) {
			return false, nil
		},
		GetRateLimitDataFunc: func(r *http.Request) RateLimitData {
			return RateLimitData{
				Limit:      100,
				Remaining:  0,
				RetryAfter: time.Second,
			}
		},
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	middleware := AdvancedMiddleware(mockLimiter, handler)
	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("expected status Too Many Requests; got %v", rr.Code)
	}
	if retryAfter := rr.Header().Get("X-RateLimit-Retry-After"); retryAfter != "1000" {
		t.Errorf("expected X-RateLimit-Retry-After header to be 1000; got %v", retryAfter)
	}
}

// Test handling errors in the rate limiter
func TestAdvancedMiddlewareHandlesError(t *testing.T) {
	mockLimiter := &MockAdvancedRateLimiter{
		IsAllowedFunc: func(r *http.Request) (bool, error) {
			return false, fmt.Errorf("rate limiter error")
		},
		GetRateLimitDataFunc: func(r *http.Request) RateLimitData {
			return RateLimitData{}
		},
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	middleware := AdvancedMiddleware(mockLimiter, handler)
	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status Internal Server Error; got %v", rr.Code)
	}
}

// Test adding rate limit headers for allowed requests
func TestAdvancedMiddlewareAddsRateLimitHeaders(t *testing.T) {
	mockLimiter := &MockAdvancedRateLimiter{
		IsAllowedFunc: func(r *http.Request) (bool, error) {
			return true, nil
		},
		GetRateLimitDataFunc: func(r *http.Request) RateLimitData {
			return RateLimitData{
				Limit:      100,
				Remaining:  99,
				RetryAfter: 0,
			}
		},
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	middleware := AdvancedMiddleware(mockLimiter, handler)
	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	if limit := rr.Header().Get("X-RateLimit-Limit"); limit != "100" {
		t.Errorf("expected X-RateLimit-Limit to be 100; got %v", limit)
	}
	if remaining := rr.Header().Get("X-RateLimit-Remaining"); remaining != "99" {
		t.Errorf("expected X-RateLimit-Remaining to be 99; got %v", remaining)
	}
}
