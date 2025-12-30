package providers

import (
	"context"
	"time"
)

// PriceData holds price information for an asset
type PriceData struct {
	Symbol           string
	Price            float64
	Currency         string
	Change24h        float64
	ChangePercent24h float64
	High24h          float64
	Low24h           float64
	Volume24h        float64
	MarketCap        float64
	UpdatedAt        time.Time
}

// HistoricalPrice holds historical price data
type HistoricalPrice struct {
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

// SearchResult holds asset search result
type SearchResult struct {
	Symbol   string
	Name     string
	Type     string
	Exchange string
	Currency string
}

// PriceProvider defines interface for price data providers
type PriceProvider interface {
	// GetPrice gets current price for a symbol
	GetPrice(ctx context.Context, symbol string) (*PriceData, error)

	// GetPrices gets prices for multiple symbols concurrently
	GetPrices(ctx context.Context, symbols []string) (map[string]*PriceData, error)

	// GetHistoricalPrices gets historical prices
	GetHistoricalPrices(ctx context.Context, symbol string, from, to time.Time) ([]HistoricalPrice, error)

	// Search searches for assets by query
	Search(ctx context.Context, query string) ([]SearchResult, error)

	// ProviderName returns the provider name
	ProviderName() string
}


