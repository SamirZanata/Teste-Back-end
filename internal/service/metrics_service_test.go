package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/back-end/quote-api/internal/domain"
	"github.com/back-end/quote-api/internal/repository"
)

func TestMetricsService_GetMetrics_InvalidLastQuotes(t *testing.T) {
	repo := &mockMetricsRepo{}
	svc := NewMetricsService(repo)

	tests := []struct {
		name   string
		param  string
	}{
		{"empty is valid (all quotes)", ""},
		{"negative", "-1"},
		{"zero", "0"},
		{"non numeric", "abc"},
		{"float", "1.5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.param == "" {
				_, err := svc.GetMetrics(context.Background(), tt.param)
				require.NoError(t, err)
				return
			}
			_, err := svc.GetMetrics(context.Background(), tt.param)
			assert.ErrorIs(t, err, ErrInvalidLastQuotes)
		})
	}
}

func TestMetricsService_GetMetrics_ValidLastQuotes(t *testing.T) {
	repo := &mockMetricsRepo{
		resp: &domain.MetricsResponse{
			ByCarrier:     []domain.CarrierMetrics{{CarrierName: "Correios", TotalQuotes: 2, TotalFreight: 41.98, AverageFreight: 20.99}},
			Cheapest:      17,
			MostExpensive: 20.99,
		},
	}
	svc := NewMetricsService(repo)

	resp, err := svc.GetMetrics(context.Background(), "5")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Len(t, resp.ByCarrier, 1)
	assert.Equal(t, 17.0, resp.Cheapest)
	assert.Equal(t, 20.99, resp.MostExpensive)
}

type mockMetricsRepo struct {
	resp *domain.MetricsResponse
}

func (m *mockMetricsRepo) CreateQuote(ctx context.Context, quote *domain.Quote) error   { return nil }
func (m *mockMetricsRepo) CreateOffer(ctx context.Context, offer *domain.QuoteOffer) error { return nil }
func (m *mockMetricsRepo) GetMetrics(ctx context.Context, lastQuotes *int) (*domain.MetricsResponse, error) {
	return m.resp, nil
}

var _ repository.QuoteRepository = (*mockMetricsRepo)(nil)
