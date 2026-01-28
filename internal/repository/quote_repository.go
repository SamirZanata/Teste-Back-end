package repository

import (
	"context"

	"github.com/back-end/quote-api/internal/domain"
)

type QuoteRepository interface {
	CreateQuote(ctx context.Context, quote *domain.Quote) error
	CreateOffer(ctx context.Context, offer *domain.QuoteOffer) error
	GetMetrics(ctx context.Context, lastQuotes *int) (*domain.MetricsResponse, error)
}
