package runhttp

import "net/http"

// HealthCheckHandler responds with a 200
type HealthCheckHandler struct {
}

// Handle responds with a 200 by default
func (h *HealthCheckHandler) Handle(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Success"))
}
