package cerberus

import (
	"net/http"
	"time"
)

// AdvancedRateLimiter is an extended version of the [RateLimiter] interface
// that provides additional data related to rate limiting. This interface
// is intended for more sophisticated rate limiting systems that need to
// expose extra information to clients, such as retry windows or usage
// limits.
//
// Implementers of this interface should provide both the core rate-limiting
// logic and a mechanism to retrieve additional rate limit details.
//
// Example use case: A system that limits client requests and also informs
// clients about their remaining quota and when they can retry.
type AdvancedRateLimiter interface {
	RateLimiter
	// GetRateLimitData retrieves additional information about the rate limit
	// for a specific request. This may include data such as remaining requests,
	// the time when the limit resets, and the total request limit.
	//
	// The returned RateLimitData can be used to enrich the response to clients,
	// providing more visibility into the rate-limiting policy (e.g., setting
	// the X-Retry-After header).
	GetRateLimitData(*http.Request) RateLimitData
}

// RateLimitData holds detailed information about the rate limiting status
// for a specific request. It is typically returned by the [AdvancedRateLimiter]
// interface to provide additional insights into the current rate limit.
//
// This struct can be used to inform clients about their remaining quota,
// the total request limit, and how long they should wait before retrying a request.
type RateLimitData struct {
	// Remaining indicates how many requests the client can still make
	// within the current rate limit window.
	Remaining int

	// Limit is the total number of requests allowed within the current
	// rate limit window.
	Limit int

	// RetryAfter specifies the amount of time a client should wait before
	// making another request. It is typically used to set the X-Retry-After
	// header in the HTTP response when rate limiting is enforced.
	RetryAfter time.Duration
}
