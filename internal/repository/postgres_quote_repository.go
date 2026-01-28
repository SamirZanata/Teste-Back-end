package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/back-end/quote-api/internal/domain"
)

type PostgresQuoteRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresQuoteRepository(pool *pgxpool.Pool) *PostgresQuoteRepository {
	return &PostgresQuoteRepository{pool: pool}
}

func (r *PostgresQuoteRepository) CreateQuote(ctx context.Context, quote *domain.Quote) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO quotes (id, zipcode, created_at) VALUES ($1, $2, NOW())`,
		quote.ID, quote.Zipcode,
	)
	return err
}

func (r *PostgresQuoteRepository) CreateOffer(ctx context.Context, offer *domain.QuoteOffer) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO quote_offers (id, quote_id, carrier_name, service, deadline_days, final_price)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		offer.ID, offer.QuoteID, offer.CarrierName, offer.Service, offer.DeadlineDays, offer.FinalPrice,
	)
	return err
}

func (r *PostgresQuoteRepository) GetMetrics(ctx context.Context, lastQuotes *int) (*domain.MetricsResponse, error) {
	limitClause := ""
	args := []interface{}{}
	argNum := 1
	if lastQuotes != nil && *lastQuotes > 0 {
		limitClause = fmt.Sprintf(" LIMIT $%d", argNum)
		args = append(args, *lastQuotes)
		argNum++
	}

	carrierQuery := fmt.Sprintf(`
		WITH selected_quotes AS (
			SELECT id FROM quotes ORDER BY created_at DESC%s
		)
		SELECT 
			o.carrier_name,
			COUNT(*)::int AS total_quotes,
			COALESCE(SUM(o.final_price), 0)::float8 AS total_freight,
			COALESCE(AVG(o.final_price), 0)::float8 AS average_freight
		FROM quote_offers o
		WHERE o.quote_id IN (SELECT id FROM selected_quotes)
		GROUP BY o.carrier_name
		ORDER BY o.carrier_name
	`, limitClause)

	rowsResult, err := r.pool.Query(ctx, carrierQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("query by carrier: %w", err)
	}
	defer rowsResult.Close()

	var byCarrier []domain.CarrierMetrics
	for rowsResult.Next() {
		var m domain.CarrierMetrics
		if err := rowsResult.Scan(&m.CarrierName, &m.TotalQuotes, &m.TotalFreight, &m.AverageFreight); err != nil {
			return nil, err
		}
		byCarrier = append(byCarrier, m)
	}
	if err := rowsResult.Err(); err != nil {
		return nil, err
	}

	minMaxQuery := fmt.Sprintf(`
		WITH selected_quotes AS (
			SELECT id FROM quotes ORDER BY created_at DESC%s
		)
		SELECT 
			COALESCE(MIN(o.final_price), 0)::float8,
			COALESCE(MAX(o.final_price), 0)::float8
		FROM quote_offers o
		WHERE o.quote_id IN (SELECT id FROM selected_quotes)
	`, limitClause)

	var cheapest, mostExpensive float64
	if err := r.pool.QueryRow(ctx, minMaxQuery, args...).Scan(&cheapest, &mostExpensive); err != nil {
		return nil, fmt.Errorf("query min/max: %w", err)
	}

	return &domain.MetricsResponse{
		ByCarrier:     byCarrier,
		Cheapest:      cheapest,
		MostExpensive: mostExpensive,
	}, nil
}

func (r *PostgresQuoteRepository) EnsureSchema(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS quotes (
			id UUID PRIMARY KEY,
			zipcode VARCHAR(20) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE TABLE IF NOT EXISTS quote_offers (
			id UUID PRIMARY KEY,
			quote_id UUID NOT NULL REFERENCES quotes(id) ON DELETE CASCADE,
			carrier_name VARCHAR(255) NOT NULL,
			service VARCHAR(255) NOT NULL,
			deadline_days INT NOT NULL,
			final_price DECIMAL(12,2) NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_quote_offers_quote_id ON quote_offers(quote_id);
		CREATE INDEX IF NOT EXISTS idx_quotes_created_at ON quotes(created_at DESC);
	`)
	return err
}
