package domain

type MetricsResponse struct {
	ByCarrier     []CarrierMetrics `json:"by_carrier"`
	Cheapest      float64          `json:"cheapest_overall"`
	MostExpensive float64          `json:"most_expensive_overall"`
}

type CarrierMetrics struct {
	CarrierName     string  `json:"carrier_name"`
	TotalQuotes     int     `json:"total_quotes"`
	TotalFreight    float64 `json:"total_freight"`
	AverageFreight  float64 `json:"average_freight"`
}
