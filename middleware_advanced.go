package cerberus

import (
	"fmt"
	"net/http"
)

// AdvancedMiddleware applies advanced rate limiting to incoming HTTP requests using an [AdvancedRateLimiter].
// In addition to enforcing rate limits, this middleware also sets custom headers in the response to provide
// information about the rate limits and retry windows.
//
// If the request exceeds the allowed rate, the middleware responds with an HTTP 429 (Too Many Requests) status code
// and sets the "X-RateLimit-Retry-After" header to indicate how long the client should wait before making
// another request. If the request is allowed, it adds headers with the remaining request quota and the total
// limit.
//
// Behavior:
//
// If the request exceeds the rate limit:
//   - Responds with an HTTP 429 (Too Many Requests) status code.
//   - Adds the "X-RateLimit-Retry-After" header, specifying the wait time before the client can retry the request.
//
// If the request is allowed:
//   - Adds the "X-RateLimit-Limit" header to indicate the total allowed requests in the current rate limit window.
//   - Adds the "X-RateLimit-Remaining" header to indicate how many requests the client can still make in the current window.
//   - Forwards the request to the next handler.
//
// If an error occurs during the rate limit check, responds with an HTTP 500 (Internal Server Error).
//
// Example usage:	http.Handle("/resource", AdvancedMiddleware(myAdvancedRateLimiter, myHandler))
func AdvancedMiddleware(rateLimiter AdvancedRateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAllowed, err := rateLimiter.IsAllowed(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		data := rateLimiter.GetRateLimitData(r)
		if !isAllowed {
			w.Header().Set("X-RateLimit-Retry-After", fmt.Sprintf("%d", data.RetryAfter.Milliseconds()))
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", data.Limit))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", data.Remaining))
		next.ServeHTTP(w, r)
	})
}
