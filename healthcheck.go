package runhttp

import "net/http"

// HealthCheckHandler responds with a 200, as required by services under Micros
type HealthCheckHandler struct {
}

// Handle responds with a 200 by default
func (h *HealthCheckHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Success"))
}
