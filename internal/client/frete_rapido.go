package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type FreteRapidoClient struct {
	baseURL       string
	token         string
	platformCode  string
	shipperCNPJ   string
	dispatcherCEP string
	httpClient    *http.Client
}

func NewFreteRapidoClient(baseURL, token, platformCode, shipperCNPJ, dispatcherCEP string) *FreteRapidoClient {
	return &FreteRapidoClient{
		baseURL:       baseURL,
		token:         token,
		platformCode:  platformCode,
		shipperCNPJ:   shipperCNPJ,
		dispatcherCEP: dispatcherCEP,
		httpClient:    &http.Client{},
	}
}

type SimulateRequest struct {
	Shipper         FRShipper         `json:"shipper"`
	Recipient       FRRecipient       `json:"recipient"`
	Dispatchers     []FRDispatcher    `json:"dispatchers"`
	SimulationType  []int            `json:"simulation_type"`
}

type FRShipper struct {
	RegisteredNumber string `json:"registered_number"`
	Token            string `json:"token"`
	PlatformCode     string `json:"platform_code"`
}

type FRRecipient struct {
	Type     int    `json:"type"`
	Country  string `json:"country"`
	Zipcode  int    `json:"zipcode"`
}

type FRDispatcher struct {
	RegisteredNumber string       `json:"registered_number"`
	Zipcode          int          `json:"zipcode"`
	Volumes          []FRVolume   `json:"volumes"`
}

type FRVolume struct {
	Amount        int     `json:"amount"`
	Category      string  `json:"category"`
	SKU           string  `json:"sku,omitempty"`
	Height        float64 `json:"height"`
	Width         float64 `json:"width"`
	Length        float64 `json:"length"`
	UnitaryPrice  float64 `json:"unitary_price"`
	UnitaryWeight float64 `json:"unitary_weight"`
}

type SimulateResponse struct {
	Dispatchers []FRDispatcherResponse `json:"dispatchers"`
}

type FRDispatcherResponse struct {
	Offers []FROffer `json:"offers"`
}

type FROffer struct {
	Carrier struct {
		Name    string `json:"name"`
		Service string `json:"service"`
	} `json:"carrier"`
	DeliveryTime struct {
		Days int `json:"days"`
	} `json:"delivery_time"`
	FinalPrice float64 `json:"final_price"`
}

func (c *FreteRapidoClient) ShipperCNPJ() string { return c.shipperCNPJ }
func (c *FreteRapidoClient) Token() string      { return c.token }
func (c *FreteRapidoClient) PlatformCode() string { return c.platformCode }
func (c *FreteRapidoClient) DispatcherCEP() string { return c.dispatcherCEP }

func (c *FreteRapidoClient) Simulate(ctx context.Context, req *SimulateRequest) (*SimulateResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := c.baseURL + "/api/v3/quote/simulate"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("frete rapido api error: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var simResp SimulateResponse
	if err := json.Unmarshal(respBody, &simResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &simResp, nil
}
