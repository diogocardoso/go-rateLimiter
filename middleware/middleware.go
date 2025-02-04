package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/diogocardoso/rate-limiter/logic"
)

func RateLimiterMiddleware(limiter *logic.RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowed, err := limiter.IsAllowed(r)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !allowed {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "you have reached the maximum number of requests or actions allowed within a certain time frame",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
