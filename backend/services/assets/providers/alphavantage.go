package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	alphaVantageBaseURL = "https://www.alphavantage.co/query"
)

// AlphaVantageClient provides stock/ETF price data from Alpha Vantage
type AlphaVantageClient struct {
	apiKey     string
	httpClient *http.Client
	isPremium  bool
}

// NewAlphaVantageClient creates a new Alpha Vantage client
func NewAlphaVantageClient(apiKey string, isPremium bool) *AlphaVantageClient {
	return &AlphaVantageClient{
		apiKey:    apiKey,
		isPremium: isPremium,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ProviderName returns the provider name
func (c *AlphaVantageClient) ProviderName() string {
	return "alphavantage"
}

// GlobalQuoteResponse represents Alpha Vantage global quote response
type GlobalQuoteResponse struct {
	GlobalQuote struct {
		Symbol           string `json:"01. symbol"`
		Open             string `json:"02. open"`
		High             string `json:"03. high"`
		Low              string `json:"04. low"`
		Price            string `json:"05. price"`
		Volume           string `json:"06. volume"`
		LatestTradingDay string `json:"07. latest trading day"`
		PreviousClose    string `json:"08. previous close"`
		Change           string `json:"09. change"`
		ChangePercent    string `json:"10. change percent"`
	} `json:"Global Quote"`
}

// GetPrice gets current price for a stock symbol
func (c *AlphaVantageClient) GetPrice(ctx context.Context, symbol string) (*PriceData, error) {
	url := fmt.Sprintf("%s?function=GLOBAL_QUOTE&symbol=%s&apikey=%s",
		alphaVantageBaseURL, symbol, c.apiKey)

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

	var result GlobalQuoteResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.GlobalQuote.Symbol == "" {
		return nil, fmt.Errorf("no data found for symbol: %s", symbol)
	}

	price, _ := strconv.ParseFloat(result.GlobalQuote.Price, 64)
	change, _ := strconv.ParseFloat(result.GlobalQuote.Change, 64)
	changePercent := 0.0
	if result.GlobalQuote.ChangePercent != "" {
		// Remove the % sign
		percentStr := result.GlobalQuote.ChangePercent
		if len(percentStr) > 0 && percentStr[len(percentStr)-1] == '%' {
			percentStr = percentStr[:len(percentStr)-1]
		}
		changePercent, _ = strconv.ParseFloat(percentStr, 64)
	}
	high, _ := strconv.ParseFloat(result.GlobalQuote.High, 64)
	low, _ := strconv.ParseFloat(result.GlobalQuote.Low, 64)
	volume, _ := strconv.ParseFloat(result.GlobalQuote.Volume, 64)

	return &PriceData{
		Symbol:           result.GlobalQuote.Symbol,
		Price:            price,
		Currency:         "USD",
		Change24h:        change,
		ChangePercent24h: changePercent,
		High24h:          high,
		Low24h:           low,
		Volume24h:        volume,
		UpdatedAt:        time.Now(),
	}, nil
}

// GetPrices gets prices for multiple symbols concurrently
func (c *AlphaVantageClient) GetPrices(ctx context.Context, symbols []string) (map[string]*PriceData, error) {
	results := make(map[string]*PriceData)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Rate limiting for free tier (5 calls/minute)
	sem := make(chan struct{}, 5)

	for _, symbol := range symbols {
		wg.Add(1)
		go func(sym string) {
			defer wg.Done()

			sem <- struct{}{} // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			price, err := c.GetPrice(ctx, sym)
			if err != nil {
				return // Skip failed symbols
			}

			mu.Lock()
			results[sym] = price
			mu.Unlock()

			// Rate limiting delay for free tier
			if !c.isPremium {
				time.Sleep(12 * time.Second) // 5 calls per minute
			}
		}(symbol)
	}

	wg.Wait()
	return results, nil
}

// TimeSeriesResponse represents Alpha Vantage time series response
type TimeSeriesResponse struct {
	MetaData struct {
		Symbol string `json:"2. Symbol"`
	} `json:"Meta Data"`
	TimeSeries map[string]struct {
		Open   string `json:"1. open"`
		High   string `json:"2. high"`
		Low    string `json:"3. low"`
		Close  string `json:"4. close"`
		Volume string `json:"5. volume"`
	} `json:"Time Series (Daily)"`
}

// GetHistoricalPrices gets historical prices for a symbol
func (c *AlphaVantageClient) GetHistoricalPrices(ctx context.Context, symbol string, from, to time.Time) ([]HistoricalPrice, error) {
	outputSize := "compact" // Last 100 days
	if to.Sub(from) > 100*24*time.Hour {
		outputSize = "full"
	}

	url := fmt.Sprintf("%s?function=TIME_SERIES_DAILY&symbol=%s&outputsize=%s&apikey=%s",
		alphaVantageBaseURL, symbol, outputSize, c.apiKey)

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

	var result TimeSeriesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var prices []HistoricalPrice
	for dateStr, data := range result.TimeSeries {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		if date.Before(from) || date.After(to) {
			continue
		}

		open, _ := strconv.ParseFloat(data.Open, 64)
		high, _ := strconv.ParseFloat(data.High, 64)
		low, _ := strconv.ParseFloat(data.Low, 64)
		closePrice, _ := strconv.ParseFloat(data.Close, 64)
		volume, _ := strconv.ParseFloat(data.Volume, 64)

		prices = append(prices, HistoricalPrice{
			Date:   date,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  closePrice,
			Volume: volume,
		})
	}

	return prices, nil
}

// SearchResponse represents Alpha Vantage search response
type SearchResponse struct {
	BestMatches []struct {
		Symbol   string `json:"1. symbol"`
		Name     string `json:"2. name"`
		Type     string `json:"3. type"`
		Region   string `json:"4. region"`
		Currency string `json:"8. currency"`
	} `json:"bestMatches"`
}

// Search searches for stocks by query
func (c *AlphaVantageClient) Search(ctx context.Context, query string) ([]SearchResult, error) {
	url := fmt.Sprintf("%s?function=SYMBOL_SEARCH&keywords=%s&apikey=%s",
		alphaVantageBaseURL, query, c.apiKey)

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

	var result SearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var results []SearchResult
	for _, match := range result.BestMatches {
		results = append(results, SearchResult{
			Symbol:   match.Symbol,
			Name:     match.Name,
			Type:     match.Type,
			Exchange: match.Region,
			Currency: match.Currency,
		})
	}

	return results, nil
}

