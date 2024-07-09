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
			if isConnectedPtr.Load() {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Repeater connected"))
			} else {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte("Repeater not connected"))
			}
		}),
	)
	logger.Info("Health check server stopped: %v", err)
}
