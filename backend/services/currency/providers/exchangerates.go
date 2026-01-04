package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	// frankfurter.app - Free, no API key required, ECB data
	frankfurterBaseURL = "https://api.frankfurter.app"
)

// ExchangeRatesClient provides exchange rate data
type ExchangeRatesClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewExchangeRatesClient creates a new exchange rates client
// apiKey parameter kept for compatibility but not used with frankfurter.app
func NewExchangeRatesClient(apiKey string) *ExchangeRatesClient {
	return &ExchangeRatesClient{
		baseURL: frankfurterBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FrankfurterResponse represents the API response
type FrankfurterResponse struct {
	Amount float64            `json:"amount"`
	Base   string             `json:"base"`
	Date   string             `json:"date"`
	Rates  map[string]float64 `json:"rates"`
}

// GetLatestRates gets the latest exchange rates for a base currency
func (c *ExchangeRatesClient) GetLatestRates(ctx context.Context, baseCurrency string) (map[string]float64, error) {
	// frankfurter.app supports: USD, EUR, GBP, JPY, etc.
	url := fmt.Sprintf("%s/latest?from=%s", c.baseURL, strings.ToUpper(baseCurrency))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result FrankfurterResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Add the base currency itself with rate 1.0
	rates := result.Rates
	if rates == nil {
		rates = make(map[string]float64)
	}
	rates[strings.ToUpper(baseCurrency)] = 1.0

	return rates, nil
}

// GetRate gets exchange rate between two currencies
func (c *ExchangeRatesClient) GetRate(ctx context.Context, from, to string) (float64, error) {
	url := fmt.Sprintf("%s/latest?from=%s&to=%s",
		c.baseURL, strings.ToUpper(from), strings.ToUpper(to))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch rate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	var result FrankfurterResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("failed to parse response: %w", err)
	}

	rate, ok := result.Rates[strings.ToUpper(to)]
	if !ok {
		return 0, fmt.Errorf("rate not found for %s to %s", from, to)
	}

	return rate, nil
}

// GetHistoricalRates gets historical rates for a specific date
func (c *ExchangeRatesClient) GetHistoricalRates(ctx context.Context, baseCurrency string, date time.Time) (map[string]float64, error) {
	dateStr := date.Format("2006-01-02")
	url := fmt.Sprintf("%s/%s?from=%s", c.baseURL, dateStr, strings.ToUpper(baseCurrency))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch historical rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result FrankfurterResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	rates := result.Rates
	if rates == nil {
		rates = make(map[string]float64)
	}
	rates[strings.ToUpper(baseCurrency)] = 1.0

	return rates, nil
}

// ConvertAmount converts an amount from one currency to another
func (c *ExchangeRatesClient) ConvertAmount(ctx context.Context, amount float64, from, to string) (float64, float64, error) {
	rate, err := c.GetRate(ctx, from, to)
	if err != nil {
		return 0, 0, err
	}

	return amount * rate, rate, nil
}
