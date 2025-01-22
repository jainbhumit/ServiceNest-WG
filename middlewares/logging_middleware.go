package middlewares

import (
	"net/http"
	"serviceNest/logger"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for preflight requests
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200") // Adjust to your frontend URL
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If the request method is OPTIONS, handle it here and return early
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		start := time.Now()
		logger.Info("Request received", map[string]interface{}{
			"method": r.Method,
			"url":    r.URL.Path,
		})

		next.ServeHTTP(w, r)

		duration := time.Since(start).String()
		logger.Info("Request processed", map[string]interface{}{
			"method":   r.Method,
			"url":      r.URL.Path,
			"duration": duration,
		})
	})
}
