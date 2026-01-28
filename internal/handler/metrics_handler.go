package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/back-end/quote-api/internal/service"
)

type MetricsHandler struct {
	svc *service.MetricsService
}

func NewMetricsHandler(svc *service.MetricsService) *MetricsHandler {
	return &MetricsHandler{svc: svc}
}

func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	lastQuotes := c.Query("last_quotes")

	resp, err := h.svc.GetMetrics(c.Request.Context(), lastQuotes)
	if err != nil {
		if err == service.ErrInvalidLastQuotes {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "O parâmetro last_quotes deve ser um número inteiro positivo (ex.: 10)",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao consultar métricas"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
