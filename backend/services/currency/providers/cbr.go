package providers

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
)

const (
	cbrBaseURL = "https://www.cbr.ru/scripts/XML_daily.asp"
)

// CBRClient provides exchange rates from Central Bank of Russia
type CBRClient struct {
	httpClient *http.Client
}

// NewCBRClient creates a new CBR client
func NewCBRClient() *CBRClient {
	return &CBRClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ValCurs represents the XML response from CBR
type ValCurs struct {
	XMLName xml.Name  `xml:"ValCurs"`
	Date    string    `xml:"Date,attr"`
	Valutes []Valute `xml:"Valute"`
}

// Valute represents a single currency in CBR response
type Valute struct {
	ID       string `xml:"ID,attr"`
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

// GetRUBRates gets exchange rates from CBR (RUB as base)
func (c *CBRClient) GetRUBRates(ctx context.Context) (map[string]float64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cbrBaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CBR rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CBR API error: %s", resp.Status)
	}

	// Decode from windows-1251 to UTF-8
	decoder := charmap.Windows1251.NewDecoder()
	reader := decoder.Reader(resp.Body)
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Remove encoding declaration since we've already decoded
	bodyStr := string(body)
	bodyStr = strings.Replace(bodyStr, `encoding="windows-1251"`, `encoding="UTF-8"`, 1)

	var valCurs ValCurs
	if err := xml.Unmarshal([]byte(bodyStr), &valCurs); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	rates := make(map[string]float64)
	rates["RUB"] = 1.0

	for _, v := range valCurs.Valutes {
		// Parse value (Russian format uses comma as decimal separator)
		valueStr := strings.Replace(v.Value, ",", ".", 1)
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			continue
		}

		// Rate is for Nominal units, so divide by nominal
		rate := value / float64(v.Nominal)
		// This gives RUB per 1 unit of foreign currency
		// We want foreign currency per 1 RUB (inverse)
		if rate > 0 {
			rates[v.CharCode] = 1 / rate
		}
	}

	return rates, nil
}

// GetUSDRUBRate gets USD/RUB rate
func (c *CBRClient) GetUSDRUBRate(ctx context.Context) (float64, error) {
	rates, err := c.GetRUBRates(ctx)
	if err != nil {
		return 0, err
	}

	usdRate, ok := rates["USD"]
	if !ok {
		return 0, fmt.Errorf("USD rate not found")
	}

	// Return RUB per 1 USD
	return 1 / usdRate, nil
}

