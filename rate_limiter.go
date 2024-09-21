package cerberus

import (
	"net/http"
)

// RateLimiter defines an interface for controlling the rate of requests
// to a resource or service. It provides a method to check if a request
// is allowed to proceed based on the implemented rate limiting logic.
type RateLimiter interface {
	// IsAllowed checks whether a request is permitted to proceed based on the
	// current rate limiting rules. It returns true if the request
	// is allowed, false otherwise. An error may be returned if
	// there are issues with the underlying rate limiting logic,
	// such as connectivity to a backend service, a cache server, or configuration errors.
	IsAllowed(*http.Request) (bool, error)
}
