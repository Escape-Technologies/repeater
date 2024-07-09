package health

import (
	"net/http"
	"sync/atomic"

	"github.com/Escape-Technologies/repeater/pkg/logger"
)

// Start the health check server
func HealthCheck(
	healthCheckPort string,
	isConnectedPtr *atomic.Bool,
) {
	logger.Info("Starting health check server on port http://localhost:%s/health", healthCheckPort)
	err := http.ListenAndServe(
		":"+healthCheckPort,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var msg string
			if isConnectedPtr.Load() {
				w.WriteHeader(http.StatusOK)
				msg = "Repeater connected"
			} else {
				w.WriteHeader(http.StatusServiceUnavailable)
				msg = "Repeater not connected"
			}
			_, err := w.Write([]byte(msg))
			if err != nil {
				logger.Debug("Error during health check: %v", err)
			}
		}),
	)
	logger.Info("Health check server stopped: %v", err)
}
