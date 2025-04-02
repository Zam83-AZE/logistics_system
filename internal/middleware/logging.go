package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Logging HTTP sorğularını loglamaq üçün middleware
func Logging(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Sorğunu emal et
			next.ServeHTTP(w, r)

			// Sorğu məlumatlarını logla
			logger.WithFields(logrus.Fields{
				"method":     r.Method,
				"path":       r.URL.Path,
				"remote_ip":  r.RemoteAddr,
				"user_agent": r.UserAgent(),
				"duration":   time.Since(start),
			}).Info("HTTP request")
		})
	}
}
