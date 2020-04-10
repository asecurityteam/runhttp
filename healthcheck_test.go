package runhttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := &HealthCheckHandler{}

	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.Handle(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}
