package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/back-end/quote-api/internal/service"
)

func TestMetricsHandler_GetMetrics_InvalidLastQuotes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/metrics?last_quotes=abc", nil)

	h := NewMetricsHandler(service.NewMetricsService(&nilQuoteRepo{}))
	h.GetMetrics(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "last_quotes")
	assert.Contains(t, w.Body.String(), "inteiro positivo")
}
