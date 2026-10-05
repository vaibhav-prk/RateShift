// Package middleware provides basic middlware chain
// auth and ratelimit
package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/vaibhav-prk/RateShift/internal/auth"
)

type contextKey string

const tenantIDKey contextKey = "tenantID"

func Auth(a *auth.Authenticator, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")

		if apiKey == "" {
			writeError(w, http.StatusUnauthorized, "missing api key")
			return
		}

		tenantID, ok := a.Authenticate(apiKey)

		if !ok {
			writeError(w, http.StatusUnauthorized, "invalid api key")
			return
		}

		ctx := context.WithValue(r.Context(), tenantIDKey, tenantID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func TenantIDFromContext(ctx context.Context) (string, bool) {
	tenantID, ok := ctx.Value(tenantIDKey).(string)
	return tenantID, ok
}
