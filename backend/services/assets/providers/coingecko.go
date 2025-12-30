package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	coinGeckoBaseURL = "https://api.coingecko.com/api/v3"
)

// CoinGeckoClient provides cryptocurrency price data from CoinGecko
type CoinGeckoClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewCoinGeckoClient creates a new CoinGecko client
func NewCoinGeckoClient(baseURL string) *CoinGeckoClient {
	if baseURL == "" {
		baseURL = coinGeckoBaseURL
	}
	return &CoinGeckoClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ProviderName returns the provider name
func (c *CoinGeckoClient) ProviderName() string {
	return "coingecko"
}

// CoinMarket represents CoinGecko market data
type CoinMarket struct {
	ID                       string  `json:"id"`
	Symbol                   string  `json:"symbol"`
	Name                     string  `json:"name"`
	CurrentPrice             float64 `json:"current_price"`
	MarketCap                float64 `json:"market_cap"`
	TotalVolume              float64 `json:"total_volume"`
	High24h                  float64 `json:"high_24h"`
	Low24h                   float64 `json:"low_24h"`
	PriceChange24h           float64 `json:"price_change_24h"`
	PriceChangePercentage24h float64 `json:"price_change_percentage_24h"`
	LastUpdated              string  `json:"last_updated"`
}

// GetPrice gets current price for a cryptocurrency
func (c *CoinGeckoClient) GetPrice(ctx context.Context, symbol string) (*PriceData, error) {
	// CoinGecko uses coin IDs, not symbols. Common mappings:
	coinID := symbolToCoinID(symbol)

	url := fmt.Sprintf("%s/coins/markets?vs_currency=usd&ids=%s", c.baseURL, coinID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch price: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var markets []CoinMarket
	if err := json.Unmarshal(body, &markets); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(markets) == 0 {
		return nil, fmt.Errorf("no data found for symbol: %s", symbol)
	}

	market := markets[0]
	return &PriceData{
		Symbol:           strings.ToUpper(market.Symbol),
		Price:            market.CurrentPrice,
		Currency:         "USD",
		Change24h:        market.PriceChange24h,
		ChangePercent24h: market.PriceChangePercentage24h,
		High24h:          market.High24h,
		Low24h:           market.Low24h,
		Volume24h:        market.TotalVolume,
		MarketCap:        market.MarketCap,
		UpdatedAt:        time.Now(),
	}, nil
}

// GetPrices gets prices for multiple cryptocurrencies concurrently
func (c *CoinGeckoClient) GetPrices(ctx context.Context, symbols []string) (map[string]*PriceData, error) {
	// Convert symbols to CoinGecko IDs
	coinIDs := make([]string, len(symbols))
	for i, sym := range symbols {
		coinIDs[i] = symbolToCoinID(sym)
	}

	url := fmt.Sprintf("%s/coins/markets?vs_currency=usd&ids=%s&per_page=250",
		c.baseURL, strings.Join(coinIDs, ","))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch prices: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var markets []CoinMarket
	if err := json.Unmarshal(body, &markets); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	results := make(map[string]*PriceData)
	for _, market := range markets {
		symbol := strings.ToUpper(market.Symbol)
		results[symbol] = &PriceData{
			Symbol:           symbol,
			Price:            market.CurrentPrice,
			Currency:         "USD",
			Change24h:        market.PriceChange24h,
			ChangePercent24h: market.PriceChangePercentage24h,
			High24h:          market.High24h,
			Low24h:           market.Low24h,
			Volume24h:        market.TotalVolume,
			MarketCap:        market.MarketCap,
			UpdatedAt:        time.Now(),
		}
	}

	return results, nil
}

// MarketChartResponse represents CoinGecko market chart response
type MarketChartResponse struct {
	Prices [][]float64 `json:"prices"`
}

// GetHistoricalPrices gets historical prices for a cryptocurrency
func (c *CoinGeckoClient) GetHistoricalPrices(ctx context.Context, symbol string, from, to time.Time) ([]HistoricalPrice, error) {
	coinID := symbolToCoinID(symbol)

	url := fmt.Sprintf("%s/coins/%s/market_chart/range?vs_currency=usd&from=%d&to=%d",
		c.baseURL, coinID, from.Unix(), to.Unix())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch historical data: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result MarketChartResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var prices []HistoricalPrice
	for _, point := range result.Prices {
		if len(point) < 2 {
			continue
		}
		timestamp := time.Unix(int64(point[0]/1000), 0)
		price := point[1]

		prices = append(prices, HistoricalPrice{
			Date:  timestamp,
			Close: price,
			Open:  price,
			High:  price,
			Low:   price,
		})
	}

	return prices, nil
}

// CoinSearchResult represents CoinGecko search result
type CoinSearchResult struct {
	Coins []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"coins"`
}

// Search searches for cryptocurrencies by query
func (c *CoinGeckoClient) Search(ctx context.Context, query string) ([]SearchResult, error) {
	url := fmt.Sprintf("%s/search?query=%s", c.baseURL, query)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result CoinSearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var results []SearchResult
	for _, coin := range result.Coins {
		results = append(results, SearchResult{
			Symbol:   strings.ToUpper(coin.Symbol),
			Name:     coin.Name,
			Type:     "crypto",
			Currency: "USD",
		})
	}

	return results, nil
}

// symbolToCoinID converts common crypto symbols to CoinGecko IDs
func symbolToCoinID(symbol string) string {
	symbol = strings.ToLower(symbol)

	// Common mappings
	mappings := map[string]string{
		"btc":   "bitcoin",
		"eth":   "ethereum",
		"usdt":  "tether",
		"bnb":   "binancecoin",
		"xrp":   "ripple",
		"usdc":  "usd-coin",
		"ada":   "cardano",
		"doge":  "dogecoin",
		"sol":   "solana",
		"dot":   "polkadot",
		"matic": "matic-network",
		"ltc":   "litecoin",
		"shib":  "shiba-inu",
		"trx":   "tron",
		"avax":  "avalanche-2",
		"link":  "chainlink",
		"atom":  "cosmos",
		"xlm":   "stellar",
		"etc":   "ethereum-classic",
		"xmr":   "monero",
	}

	if id, ok := mappings[symbol]; ok {
		return id
	}
	return symbol
}

// PriceManager manages multiple price providers
type PriceManager struct {
	stockProvider  PriceProvider
	cryptoProvider PriceProvider
	mu             sync.RWMutex
}

// NewPriceManager creates a new price manager
func NewPriceManager(stockProvider, cryptoProvider PriceProvider) *PriceManager {
	return &PriceManager{
		stockProvider:  stockProvider,
		cryptoProvider: cryptoProvider,
	}
}

// GetPrice gets price from appropriate provider based on asset type
func (m *PriceManager) GetPrice(ctx context.Context, symbol string, assetType string) (*PriceData, error) {
	switch assetType {
	case "crypto":
		return m.cryptoProvider.GetPrice(ctx, symbol)
	case "stock", "etf":
		return m.stockProvider.GetPrice(ctx, symbol)
	default:
		// Try stock first, then crypto
		price, err := m.stockProvider.GetPrice(ctx, symbol)
		if err == nil {
			return price, nil
		}
		return m.cryptoProvider.GetPrice(ctx, symbol)
	}
}

// GetPricesConcurrently gets prices for multiple assets concurrently
func (m *PriceManager) GetPricesConcurrently(ctx context.Context, assets map[string]string) (map[string]*PriceData, error) {
	results := make(map[string]*PriceData)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for symbol, assetType := range assets {
		wg.Add(1)
		go func(sym, aType string) {
			defer wg.Done()

			price, err := m.GetPrice(ctx, sym, aType)
			if err != nil {
				return
			}

			mu.Lock()
			results[sym] = price
			mu.Unlock()
		}(symbol, assetType)
	}

	wg.Wait()
	return results, nil
}

// GetStockProvider returns the stock provider
func (m *PriceManager) GetStockProvider() PriceProvider {
	return m.stockProvider
}

// GetCryptoProvider returns the crypto provider
func (m *PriceManager) GetCryptoProvider() PriceProvider {
	return m.cryptoProvider
}

