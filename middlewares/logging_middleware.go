package middlewares

import (
	"net/http"
	"serviceNest/logger"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
