package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/radmickey/money-control/backend/pkg/cache"
	"github.com/radmickey/money-control/backend/services/currency/models"
	"github.com/radmickey/money-control/backend/services/currency/providers"
	"github.com/radmickey/money-control/backend/services/currency/repository"
)

const (
	rateCacheTTL = 1 * time.Hour
)

// CurrencyService handles currency business logic
type CurrencyService struct {
	currencyRepo     *repository.CurrencyRepository
	rateRepo         *repository.ExchangeRateRepository
	historyRepo      *repository.RateHistoryRepository
	ratesClient      *providers.ExchangeRatesClient
	cache            *cache.Cache
	updateTicker     *time.Ticker
	stopChan         chan struct{}
	defaultBase      string
	mu               sync.RWMutex
}

// NewCurrencyService creates a new currency service
func NewCurrencyService(
	currencyRepo *repository.CurrencyRepository,
	rateRepo *repository.ExchangeRateRepository,
	historyRepo *repository.RateHistoryRepository,
	ratesClient *providers.ExchangeRatesClient,
	redisCache *cache.Cache,
	defaultBase string,
) *CurrencyService {
	if defaultBase == "" {
		defaultBase = "USD"
	}

	s := &CurrencyService{
		currencyRepo: currencyRepo,
		rateRepo:     rateRepo,
		historyRepo:  historyRepo,
		ratesClient:  ratesClient,
		cache:        redisCache,
		defaultBase:  defaultBase,
		stopChan:     make(chan struct{}),
	}

	return s
}

// StartRateUpdater starts the background rate updater
func (s *CurrencyService) StartRateUpdater(interval time.Duration) {
	s.updateTicker = time.NewTicker(interval)

	go func() {
		// Initial update
		if _, err := s.RefreshRates(context.Background(), s.defaultBase); err != nil {
			log.Printf("Failed to refresh rates: %v", err)
		}

		for {
			select {
			case <-s.updateTicker.C:
				if _, err := s.RefreshRates(context.Background(), s.defaultBase); err != nil {
					log.Printf("Failed to refresh rates: %v", err)
				}
			case <-s.stopChan:
				s.updateTicker.Stop()
				return
			}
		}
	}()
}

// Stop stops the rate updater
func (s *CurrencyService) Stop() {
	close(s.stopChan)
}

// GetExchangeRate gets exchange rate between two currencies
func (s *CurrencyService) GetExchangeRate(ctx context.Context, from, to string) (*models.ExchangeRate, error) {
	// Check cache first
	cacheKey := cache.ExchangeRateKey(from, to)
	var cachedRate models.ExchangeRate
	if err := s.cache.Get(ctx, cacheKey, &cachedRate); err == nil {
		return &cachedRate, nil
	}

	// Check database
	rate, err := s.rateRepo.GetRate(ctx, from, to)
	if err == nil {
		// Cache the result
		_ = s.cache.Set(ctx, cacheKey, rate, rateCacheTTL)
		return rate, nil
	}

	// Try to calculate from available rates
	if from != s.defaultBase && to != s.defaultBase {
		fromRate, err := s.rateRepo.GetRate(ctx, s.defaultBase, from)
		if err != nil {
			return nil, err
		}
		toRate, err := s.rateRepo.GetRate(ctx, s.defaultBase, to)
		if err != nil {
			return nil, err
		}

		calculatedRate := toRate.Rate / fromRate.Rate
		result := &models.ExchangeRate{
			FromCurrency: from,
			ToCurrency:   to,
			Rate:         calculatedRate,
			UpdatedAt:    time.Now(),
		}
		_ = s.cache.Set(ctx, cacheKey, result, rateCacheTTL)
		return result, nil
	}

	return nil, err
}

// GetMultipleExchangeRates gets exchange rates for multiple currencies
func (s *CurrencyService) GetMultipleExchangeRates(ctx context.Context, baseCurrency string, targetCurrencies []string) (map[string]float64, time.Time, error) {
	// Check cache
	cacheKey := cache.ExchangeRatesKey(baseCurrency)
	var cachedRates map[string]float64
	if err := s.cache.Get(ctx, cacheKey, &cachedRates); err == nil {
		// Filter to requested currencies
		result := make(map[string]float64)
		for _, currency := range targetCurrencies {
			if rate, ok := cachedRates[currency]; ok {
				result[currency] = rate
			}
		}
		return result, time.Now(), nil
	}

	// Get from database
	allRates, lastUpdate, err := s.rateRepo.GetRatesForBase(ctx, baseCurrency)
	if err != nil {
		return nil, time.Time{}, err
	}

	// Cache all rates
	_ = s.cache.Set(ctx, cacheKey, allRates, rateCacheTTL)

	// Filter to requested currencies
	if len(targetCurrencies) > 0 {
		result := make(map[string]float64)
		for _, currency := range targetCurrencies {
			if rate, ok := allRates[currency]; ok {
				result[currency] = rate
			}
		}
		return result, lastUpdate, nil
	}

	return allRates, lastUpdate, nil
}

// ConvertAmount converts an amount from one currency to another
func (s *CurrencyService) ConvertAmount(ctx context.Context, amount float64, from, to string) (float64, float64, error) {
	if from == to {
		return amount, 1.0, nil
	}

	rate, err := s.GetExchangeRate(ctx, from, to)
	if err != nil {
		return 0, 0, err
	}

	return amount * rate.Rate, rate.Rate, nil
}

// ConvertMultipleAmounts converts multiple amounts to a target currency
func (s *CurrencyService) ConvertMultipleAmounts(ctx context.Context, amounts []AmountToConvert, toCurrency string) ([]ConvertedAmount, float64, error) {
	results := make([]ConvertedAmount, len(amounts))
	var total float64

	for i, amount := range amounts {
		converted, rate, err := s.ConvertAmount(ctx, amount.Amount, amount.FromCurrency, toCurrency)
		if err != nil {
			return nil, 0, err
		}

		results[i] = ConvertedAmount{
			ID:              amount.ID,
			OriginalAmount:  amount.Amount,
			FromCurrency:    amount.FromCurrency,
			ConvertedAmount: converted,
			RateUsed:        rate,
		}
		total += converted
	}

	return results, total, nil
}

// ListSupportedCurrencies lists all supported currencies
func (s *CurrencyService) ListSupportedCurrencies(ctx context.Context, includeCrypto bool) ([]models.Currency, error) {
	return s.currencyRepo.GetAllCurrencies(ctx, includeCrypto)
}

// GetRateHistory gets historical rates for a currency pair
func (s *CurrencyService) GetRateHistory(ctx context.Context, from, to string, startDate, endDate time.Time) ([]models.RateHistory, error) {
	return s.historyRepo.GetHistory(ctx, from, to, startDate, endDate)
}

// RefreshRates refreshes exchange rates from the API
func (s *CurrencyService) RefreshRates(ctx context.Context, baseCurrency string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rates, err := s.ratesClient.GetLatestRates(ctx, baseCurrency)
	if err != nil {
		return 0, err
	}

	// Save to database
	if err := s.rateRepo.UpsertRates(ctx, baseCurrency, rates); err != nil {
		return 0, err
	}

	// Invalidate cache
	cacheKey := cache.ExchangeRatesKey(baseCurrency)
	_ = s.cache.Delete(ctx, cacheKey)

	// Cache new rates
	_ = s.cache.Set(ctx, cacheKey, rates, rateCacheTTL)

	log.Printf("Updated %d exchange rates for base %s", len(rates), baseCurrency)

	return len(rates), nil
}

// GetLastUpdateTime gets the last update time for rates
func (s *CurrencyService) GetLastUpdateTime(ctx context.Context, baseCurrency string) (time.Time, error) {
	return s.rateRepo.GetLastUpdateTime(ctx, baseCurrency)
}

// AmountToConvert holds an amount to convert
type AmountToConvert struct {
	ID           string
	Amount       float64
	FromCurrency string
}

// ConvertedAmount holds a converted amount
type ConvertedAmount struct {
	ID              string
	OriginalAmount  float64
	FromCurrency    string
	ConvertedAmount float64
	RateUsed        float64
}

