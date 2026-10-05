package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vaibhav-prk/RateShift/internal/limiter"
)

func RateLimit(l limiter.RateLimiter, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID, ok := TenantIDFromContext(r.Context())

		if !ok {
			writeError(w, http.StatusUnauthorized, "tenant not found in context")
			return
		}

		if l == nil {
			log.Printf("rate limiting would happen for tenant: %s", tenantID)
			next.ServeHTTP(w, r)
			return
		}

		decision, err := l.Allow(r.Context(), tenantID)
		if err != nil {
			log.Printf("rate limiter error for tenant %s: %v", tenantID, err)
			next.ServeHTTP(w, r)
			return
		}

		if !decision.Allowed {
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("Retry-After", fmt.Sprintf("%d", int(time.Until(decision.ResetAfter).Seconds())))
			writeError(w, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}

		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", decision.Remaining))
		next.ServeHTTP(w, r)
	}
}
