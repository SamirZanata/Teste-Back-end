package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/back-end/quote-api/internal/domain"
	"github.com/back-end/quote-api/internal/repository"
)

var ErrInvalidLastQuotes = errors.New("last_quotes deve ser um número inteiro positivo")

type MetricsService struct {
	repo repository.QuoteRepository
}

func NewMetricsService(repo repository.QuoteRepository) *MetricsService {
	return &MetricsService{repo: repo}
}

func (s *MetricsService) GetMetrics(ctx context.Context, lastQuotesRaw string) (*domain.MetricsResponse, error) {
	var lastQuotes *int
	if lastQuotesRaw != "" {
		n, err := strconv.Atoi(lastQuotesRaw)
		if err != nil || n < 1 {
			return nil, ErrInvalidLastQuotes
		}
		lastQuotes = &n
	}
	return s.repo.GetMetrics(ctx, lastQuotes)
}
