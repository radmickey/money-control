package repository

import (
	"context"
	"errors"
	"time"

	"github.com/radmickey/money-control/backend/services/currency/models"
	"gorm.io/gorm"
)

var (
	ErrRateNotFound = errors.New("exchange rate not found")
)

// CurrencyRepository handles database operations for currencies
type CurrencyRepository struct {
	db *gorm.DB
}

// NewCurrencyRepository creates a new currency repository
func NewCurrencyRepository(db *gorm.DB) *CurrencyRepository {
	return &CurrencyRepository{db: db}
}

// SeedCurrencies seeds the currencies table
func (r *CurrencyRepository) SeedCurrencies(ctx context.Context) error {
	for _, currency := range models.SupportedCurrencies {
		if err := r.db.WithContext(ctx).
			Where("code = ?", currency.Code).
			FirstOrCreate(&currency).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetAllCurrencies gets all supported currencies
func (r *CurrencyRepository) GetAllCurrencies(ctx context.Context, includeCrypto bool) ([]models.Currency, error) {
	var currencies []models.Currency
	query := r.db.WithContext(ctx)
	if !includeCrypto {
		query = query.Where("is_crypto = false")
	}
	if err := query.Order("code").Find(&currencies).Error; err != nil {
		return nil, err
	}
	return currencies, nil
}

// ExchangeRateRepository handles exchange rate operations
type ExchangeRateRepository struct {
	db *gorm.DB
}

// NewExchangeRateRepository creates a new exchange rate repository
func NewExchangeRateRepository(db *gorm.DB) *ExchangeRateRepository {
	return &ExchangeRateRepository{db: db}
}

// UpsertRate inserts or updates an exchange rate
func (r *ExchangeRateRepository) UpsertRate(ctx context.Context, rate *models.ExchangeRate) error {
	return r.db.WithContext(ctx).
		Where("from_currency = ? AND to_currency = ?", rate.FromCurrency, rate.ToCurrency).
		Assign(rate).
		FirstOrCreate(rate).Error
}

// UpsertRates batch upserts multiple exchange rates
func (r *ExchangeRateRepository) UpsertRates(ctx context.Context, baseCurrency string, rates map[string]float64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		for currency, rate := range rates {
			exchangeRate := &models.ExchangeRate{
				FromCurrency: baseCurrency,
				ToCurrency:   currency,
				Rate:         rate,
				UpdatedAt:    now,
			}
			if err := tx.Where("from_currency = ? AND to_currency = ?", baseCurrency, currency).
				Assign(exchangeRate).
				FirstOrCreate(exchangeRate).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetRate gets an exchange rate
func (r *ExchangeRateRepository) GetRate(ctx context.Context, from, to string) (*models.ExchangeRate, error) {
	var rate models.ExchangeRate
	if err := r.db.WithContext(ctx).
		Where("from_currency = ? AND to_currency = ?", from, to).
		First(&rate).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRateNotFound
		}
		return nil, err
	}
	return &rate, nil
}

// GetRatesForBase gets all rates for a base currency
func (r *ExchangeRateRepository) GetRatesForBase(ctx context.Context, baseCurrency string) (map[string]float64, time.Time, error) {
	var rates []models.ExchangeRate
	if err := r.db.WithContext(ctx).
		Where("from_currency = ?", baseCurrency).
		Find(&rates).Error; err != nil {
		return nil, time.Time{}, err
	}

	result := make(map[string]float64)
	var lastUpdate time.Time
	for _, rate := range rates {
		result[rate.ToCurrency] = rate.Rate
		if rate.UpdatedAt.After(lastUpdate) {
			lastUpdate = rate.UpdatedAt
		}
	}
	return result, lastUpdate, nil
}

// GetLastUpdateTime gets the last update time for rates
func (r *ExchangeRateRepository) GetLastUpdateTime(ctx context.Context, baseCurrency string) (time.Time, error) {
	var rate models.ExchangeRate
	if err := r.db.WithContext(ctx).
		Where("from_currency = ?", baseCurrency).
		Order("updated_at DESC").
		First(&rate).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}
	return rate.UpdatedAt, nil
}

// RateHistoryRepository handles rate history operations
type RateHistoryRepository struct {
	db *gorm.DB
}

// NewRateHistoryRepository creates a new rate history repository
func NewRateHistoryRepository(db *gorm.DB) *RateHistoryRepository {
	return &RateHistoryRepository{db: db}
}

// Create creates a rate history entry
func (r *RateHistoryRepository) Create(ctx context.Context, history *models.RateHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// GetHistory gets rate history for a currency pair
func (r *RateHistoryRepository) GetHistory(ctx context.Context, from, to string, startDate, endDate time.Time) ([]models.RateHistory, error) {
	var history []models.RateHistory
	query := r.db.WithContext(ctx).
		Where("from_currency = ? AND to_currency = ?", from, to)

	if !startDate.IsZero() {
		query = query.Where("date >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("date <= ?", endDate)
	}

	if err := query.Order("date DESC").Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

