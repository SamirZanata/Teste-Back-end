package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/back-end/quote-api/internal/client"
	"github.com/back-end/quote-api/internal/domain"
	"github.com/back-end/quote-api/internal/repository"
)

type QuoteService struct {
	repo   repository.QuoteRepository
	client *client.FreteRapidoClient
}

func NewQuoteService(repo repository.QuoteRepository, frClient *client.FreteRapidoClient) *QuoteService {
	return &QuoteService{repo: repo, client: frClient}
}

func (s *QuoteService) CreateQuote(ctx context.Context, req *domain.QuoteRequest) (*domain.QuoteResponse, error) {
	if err := s.validateZipcode(req.Recipient.Address.Zipcode); err != nil {
		return nil, err
	}

	recipientZipcode, err := zipcodeToInt(req.Recipient.Address.Zipcode)
	if err != nil {
		return nil, fmt.Errorf("zipcode inválido: deve conter apenas 8 dígitos numéricos")
	}

	frReq := s.buildFreteRapidoRequest(recipientZipcode, req)
	simResp, err := s.client.Simulate(ctx, frReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter cotação do Frete Rápido: %w", err)
	}

	offers := s.extractOffers(simResp)
	if len(offers) == 0 {
		return &domain.QuoteResponse{Carrier: []domain.CarrierOffer{}}, nil
	}

	quoteID := uuid.New()
	quote := &domain.Quote{ID: quoteID, Zipcode: req.Recipient.Address.Zipcode}
	if err := s.repo.CreateQuote(ctx, quote); err != nil {
		return nil, fmt.Errorf("erro ao salvar cotação: %w", err)
	}

	for _, o := range offers {
		offer := &domain.QuoteOffer{
			ID:           uuid.New(),
			QuoteID:      quoteID,
			CarrierName:  o.Name,
			Service:      o.Service,
			DeadlineDays: parseIntDeadline(o.Deadline),
			FinalPrice:   o.Price,
		}
		if err := s.repo.CreateOffer(ctx, offer); err != nil {
			return nil, fmt.Errorf("erro ao salvar oferta: %w", err)
		}
	}

	return &domain.QuoteResponse{Carrier: offers}, nil
}

func (s *QuoteService) validateZipcode(zipcode string) error {
	if len(zipcode) != 8 {
		return fmt.Errorf("zipcode deve ter exatamente 8 caracteres")
	}
	for _, c := range zipcode {
		if c < '0' || c > '9' {
			return fmt.Errorf("zipcode deve conter apenas dígitos numéricos")
		}
	}
	return nil
}

func zipcodeToInt(z string) (int, error) {
	i, err := strconv.Atoi(z)
	if err != nil || len(z) != 8 {
		return 0, fmt.Errorf("zipcode inválido")
	}
	return i, nil
}

func parseIntDeadline(deadline string) int {
	d, _ := strconv.Atoi(deadline)
	if d < 0 {
		return 0
	}
	return d
}

func (s *QuoteService) buildFreteRapidoRequest(recipientZipcode int, req *domain.QuoteRequest) *client.SimulateRequest {
	volumes := make([]client.FRVolume, len(req.Volumes))
	for i, v := range req.Volumes {
		volumes[i] = client.FRVolume{
			Amount:        v.Amount,
			Category:      strconv.Itoa(v.Category),
			SKU:           v.SKU,
			Height:        v.Height,
			Width:         v.Width,
			Length:        v.Length,
			UnitaryPrice:  v.Price,
			UnitaryWeight: v.UnitaryWeight,
		}
	}
	dispatcherZipcode, _ := strconv.Atoi(s.client.DispatcherCEP())
	if dispatcherZipcode == 0 {
		dispatcherZipcode = 29161376
	}
	return &client.SimulateRequest{
		Shipper: client.FRShipper{
			RegisteredNumber: s.client.ShipperCNPJ(),
			Token:            s.client.Token(),
			PlatformCode:     s.client.PlatformCode(),
		},
		Recipient: client.FRRecipient{
			Type:    0,
			Country: "BRA",
			Zipcode: recipientZipcode,
		},
		Dispatchers: []client.FRDispatcher{
			{
				RegisteredNumber: s.client.ShipperCNPJ(),
				Zipcode:          dispatcherZipcode,
				Volumes:          volumes,
			},
		},
		SimulationType: []int{0},
	}
}

func (s *QuoteService) extractOffers(resp *client.SimulateResponse) []domain.CarrierOffer {
	var out []domain.CarrierOffer
	for _, d := range resp.Dispatchers {
		for _, o := range d.Offers {
			days := 0
			if o.DeliveryTime.Days > 0 {
				days = o.DeliveryTime.Days
			}
			out = append(out, domain.CarrierOffer{
				Name:     o.Carrier.Name,
				Service:  o.Carrier.Service,
				Deadline: strconv.Itoa(days),
				Price:    o.FinalPrice,
			})
		}
	}
	return out
}
