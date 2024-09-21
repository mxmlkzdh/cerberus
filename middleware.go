package cerberus

import "net/http"

// Middleware applies rate limiting to incoming HTTP requests using the provided [RateLimiter].
// It checks if the request is allowed to proceed based on the rate limiting rules defined by
// the rateLimiter. If the request exceeds the allowed rate, it responds with an HTTP 429 (Too Many Requests)
// status code. If an error occurs while checking the rate limit,
// it responds with an HTTP 500 (Internal Server Error).
//
// Behavior:
//   - If the request is allowed by the rate limiter, it is forwarded to the next handler in the chain.
//   - If the request exceeds the rate limit, an HTTP 429 (Too Many Requests) response is returned.
//   - If the rate limiter encounters an error, an HTTP 500 (Internal Server Error) response is returned.
//
// Example usage: http.Handle("/resource", cerberus.Middleware(myRateLimiter, myHandler))
func Middleware(rateLimiter RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAllowed, err := rateLimiter.IsAllowed(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !isAllowed {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
