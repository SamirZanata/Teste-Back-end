package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/back-end/quote-api/internal/domain"
	"github.com/back-end/quote-api/internal/repository"
	"github.com/back-end/quote-api/internal/service"
)

func TestQuoteHandler_CreateQuote_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/quote", bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")

	h := NewQuoteHandler(service.NewQuoteService(&nilQuoteRepo{}, nil))
	h.CreateQuote(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "inválidos")
}

func TestQuoteHandler_CreateQuote_ValidationError_MissingZipcode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := `{"recipient":{"address":{}},"volumes":[{"category":7,"amount":1,"unitary_weight":5,"price":349,"height":0.2,"width":0.2,"length":0.2}]}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/quote", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h := NewQuoteHandler(service.NewQuoteService(&nilQuoteRepo{}, nil))
	h.CreateQuote(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "zipcode")
}

type nilQuoteRepo struct{}

func (n *nilQuoteRepo) CreateQuote(ctx context.Context, quote *domain.Quote) error { return nil }
func (n *nilQuoteRepo) CreateOffer(ctx context.Context, offer *domain.QuoteOffer) error { return nil }
func (n *nilQuoteRepo) GetMetrics(ctx context.Context, lastQuotes *int) (*domain.MetricsResponse, error) {
	return &domain.MetricsResponse{}, nil
}

var _ repository.QuoteRepository = (*nilQuoteRepo)(nil)
