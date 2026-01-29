package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/back-end/quote-api/internal/client"
	"github.com/back-end/quote-api/internal/domain"
	"github.com/back-end/quote-api/internal/repository"
)

func TestQuoteService_CreateQuote_ValidZipcode(t *testing.T) {
	offers := []map[string]interface{}{
		{"carrier": map[string]string{"name": "EXPRESSO FR", "service": "Rodoviário"}, "delivery_time": map[string]int{"days": 3}, "final_price": 17.0},
		{"carrier": map[string]string{"name": "Correios", "service": "SEDEX"}, "delivery_time": map[string]int{"days": 1}, "final_price": 20.99},
	}
	respBody := map[string]interface{}{
		"dispatchers": []map[string]interface{}{
			{"offers": offers},
		},
	}
	body, _ := json.Marshal(respBody)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}))
	defer server.Close()

	repo := &mockQuoteRepo{}
	frClient := client.NewFreteRapidoClient(server.URL, "token", "code", "25438296000158", "29161376")
	svc := NewQuoteService(repo, frClient)

	req := &domain.QuoteRequest{
		Recipient: domain.QuoteRecipient{
			Address: domain.QuoteAddress{Zipcode: "01311000"},
		},
		Volumes: []domain.QuoteVolume{
			{Category: 7, Amount: 1, UnitaryWeight: 5, Price: 349, SKU: "abc", Height: 0.2, Width: 0.2, Length: 0.2},
		},
	}

	resp, err := svc.CreateQuote(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Len(t, resp.Carrier, 2)
	assert.Equal(t, "EXPRESSO FR", resp.Carrier[0].Name)
	assert.Equal(t, "Rodoviário", resp.Carrier[0].Service)
	assert.Equal(t, "3", resp.Carrier[0].Deadline)
	assert.Equal(t, 17.0, resp.Carrier[0].Price)
	assert.Equal(t, "Correios", resp.Carrier[1].Name)
	assert.Equal(t, "SEDEX", resp.Carrier[1].Service)
	assert.Equal(t, "1", resp.Carrier[1].Deadline)
	assert.Equal(t, 20.99, resp.Carrier[1].Price)
	assert.Equal(t, 1, repo.createQuoteCalls)
	assert.Equal(t, 2, repo.createOfferCalls)
}

func TestQuoteService_CreateQuote_InvalidZipcode_Length(t *testing.T) {
	repo := &mockQuoteRepo{}
	frClient := client.NewFreteRapidoClient("http://localhost", "t", "c", "25438296000158", "29161376")
	svc := NewQuoteService(repo, frClient)

	req := &domain.QuoteRequest{
		Recipient: domain.QuoteRecipient{Address: domain.QuoteAddress{Zipcode: "01311"}},
		Volumes:   []domain.QuoteVolume{{Category: 7, Amount: 1, UnitaryWeight: 5, Price: 349, Height: 0.2, Width: 0.2, Length: 0.2}},
	}

	resp, err := svc.CreateQuote(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "zipcode")
	assert.Zero(t, repo.createQuoteCalls)
}

func TestQuoteService_CreateQuote_InvalidZipcode_NonNumeric(t *testing.T) {
	repo := &mockQuoteRepo{}
	frClient := client.NewFreteRapidoClient("http://localhost", "t", "c", "25438296000158", "29161376")
	svc := NewQuoteService(repo, frClient)

	req := &domain.QuoteRequest{
		Recipient: domain.QuoteRecipient{Address: domain.QuoteAddress{Zipcode: "01311abc"}},
		Volumes:   []domain.QuoteVolume{{Category: 7, Amount: 1, UnitaryWeight: 5, Price: 349, Height: 0.2, Width: 0.2, Length: 0.2}},
	}

	resp, err := svc.CreateQuote(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "zipcode")
	assert.Zero(t, repo.createQuoteCalls)
}

type mockQuoteRepo struct {
	createQuoteCalls int
	createOfferCalls int
}

func (m *mockQuoteRepo) CreateQuote(ctx context.Context, quote *domain.Quote) error {
	m.createQuoteCalls++
	return nil
}
func (m *mockQuoteRepo) CreateOffer(ctx context.Context, offer *domain.QuoteOffer) error {
	m.createOfferCalls++
	return nil
}
func (m *mockQuoteRepo) GetMetrics(ctx context.Context, lastQuotes *int) (*domain.MetricsResponse, error) {
	return nil, nil
}

var _ repository.QuoteRepository = (*mockQuoteRepo)(nil)
