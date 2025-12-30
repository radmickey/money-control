package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	exchangeRatesBaseURL = "https://api.exchangeratesapi.io/v1"
)

// ExchangeRatesClient provides exchange rate data
type ExchangeRatesClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// NewExchangeRatesClient creates a new exchange rates client
func NewExchangeRatesClient(apiKey string) *ExchangeRatesClient {
	return &ExchangeRatesClient{
		apiKey:  apiKey,
		baseURL: exchangeRatesBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// LatestRatesResponse represents the API response for latest rates
type LatestRatesResponse struct {
	Success   bool               `json:"success"`
	Timestamp int64              `json:"timestamp"`
	Base      string             `json:"base"`
	Date      string             `json:"date"`
	Rates     map[string]float64 `json:"rates"`
	Error     *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// GetLatestRates gets the latest exchange rates for a base currency
func (c *ExchangeRatesClient) GetLatestRates(ctx context.Context, baseCurrency string) (map[string]float64, error) {
	url := fmt.Sprintf("%s/latest?access_key=%s&base=%s",
		c.baseURL, c.apiKey, baseCurrency)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result LatestRatesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Free tier only supports EUR as base
	// If different base requested, we need to convert
	if result.Base != baseCurrency && len(result.Rates) > 0 {
		return c.convertToBase(result.Rates, result.Base, baseCurrency), nil
	}

	return result.Rates, nil
}

// GetRate gets exchange rate between two currencies
func (c *ExchangeRatesClient) GetRate(ctx context.Context, from, to string) (float64, error) {
	rates, err := c.GetLatestRates(ctx, from)
	if err != nil {
		return 0, err
	}

	rate, ok := rates[to]
	if !ok {
		return 0, fmt.Errorf("rate not found for %s to %s", from, to)
	}

	return rate, nil
}

// HistoricalRatesResponse represents historical rates response
type HistoricalRatesResponse struct {
	Success    bool               `json:"success"`
	Historical bool               `json:"historical"`
	Date       string             `json:"date"`
	Timestamp  int64              `json:"timestamp"`
	Base       string             `json:"base"`
	Rates      map[string]float64 `json:"rates"`
}

// GetHistoricalRates gets historical rates for a specific date
func (c *ExchangeRatesClient) GetHistoricalRates(ctx context.Context, baseCurrency string, date time.Time) (map[string]float64, error) {
	dateStr := date.Format("2006-01-02")
	url := fmt.Sprintf("%s/%s?access_key=%s&base=%s",
		c.baseURL, dateStr, c.apiKey, baseCurrency)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch historical rates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result HistoricalRatesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Base != baseCurrency && len(result.Rates) > 0 {
		return c.convertToBase(result.Rates, result.Base, baseCurrency), nil
	}

	return result.Rates, nil
}

// convertToBase converts rates from one base currency to another
func (c *ExchangeRatesClient) convertToBase(rates map[string]float64, currentBase, newBase string) map[string]float64 {
	// Get the rate of the new base in terms of current base
	newBaseRate, ok := rates[newBase]
	if !ok || newBaseRate == 0 {
		return rates
	}

	converted := make(map[string]float64)
	for currency, rate := range rates {
		if currency == newBase {
			converted[currentBase] = 1 / newBaseRate
		} else {
			converted[currency] = rate / newBaseRate
		}
	}
	converted[currentBase] = 1 / newBaseRate

	return converted
}

// ConvertAmount converts an amount from one currency to another
func (c *ExchangeRatesClient) ConvertAmount(ctx context.Context, amount float64, from, to string) (float64, float64, error) {
	rate, err := c.GetRate(ctx, from, to)
	if err != nil {
		return 0, 0, err
	}

	return amount * rate, rate, nil
}

