package domain

import "github.com/google/uuid"

type QuoteRequest struct {
	Recipient QuoteRecipient `json:"recipient" binding:"required"`
	Volumes   []QuoteVolume  `json:"volumes" binding:"required,min=1,dive"`
}

type QuoteRecipient struct {
	Address QuoteAddress `json:"address" binding:"required"`
}

type QuoteAddress struct {
	Zipcode string `json:"zipcode" binding:"required,len=8"`
}

type QuoteVolume struct {
	Category       int     `json:"category" binding:"required,min=1"`
	Amount         int     `json:"amount" binding:"required,min=1"`
	UnitaryWeight  float64 `json:"unitary_weight" binding:"required,gt=0"`
	Price          float64 `json:"price" binding:"required,gte=0"`
	SKU            string  `json:"sku" binding:"omitempty"`
	Height         float64 `json:"height" binding:"required,gt=0"`
	Width          float64 `json:"width" binding:"required,gt=0"`
	Length         float64 `json:"length" binding:"required,gt=0"`
}

type CarrierOffer struct {
	Name     string  `json:"name"`
	Service  string  `json:"service"`
	Deadline string  `json:"deadline"`
	Price    float64 `json:"price"`
}

type QuoteResponse struct {
	Carrier []CarrierOffer `json:"carrier"`
}

type Quote struct {
	ID        uuid.UUID
	Zipcode   string
}

type QuoteOffer struct {
	ID           uuid.UUID
	QuoteID      uuid.UUID
	CarrierName  string
	Service      string
	DeadlineDays int
	FinalPrice   float64
}
