package models

import (
	"time"
)

// ExchangeRate represents an exchange rate between two currencies
type ExchangeRate struct {
	ID           string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FromCurrency string    `gorm:"size:3;not null;index" json:"from_currency"`
	ToCurrency   string    `gorm:"size:3;not null;index" json:"to_currency"`
	Rate         float64   `gorm:"type:decimal(20,10);not null" json:"rate"`
	UpdatedAt    time.Time `gorm:"index" json:"updated_at"`
}

// TableName returns the table name for GORM
func (ExchangeRate) TableName() string {
	return "exchange_rates"
}

// RateHistory stores historical exchange rates
type RateHistory struct {
	ID           string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FromCurrency string    `gorm:"size:3;not null;index" json:"from_currency"`
	ToCurrency   string    `gorm:"size:3;not null;index" json:"to_currency"`
	Rate         float64   `gorm:"type:decimal(20,10);not null" json:"rate"`
	Date         time.Time `gorm:"type:date;not null;index" json:"date"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName returns the table name for GORM
func (RateHistory) TableName() string {
	return "rate_history"
}

// Currency represents a supported currency
type Currency struct {
	Code          string `gorm:"size:3;primary_key" json:"code"`
	Name          string `gorm:"size:100;not null" json:"name"`
	Symbol        string `gorm:"size:10" json:"symbol"`
	DecimalPlaces int    `gorm:"default:2" json:"decimal_places"`
	IsCrypto      bool   `gorm:"default:false" json:"is_crypto"`
}

// TableName returns the table name for GORM
func (Currency) TableName() string {
	return "currencies"
}

// SupportedCurrencies is a list of commonly supported currencies
var SupportedCurrencies = []Currency{
	{Code: "USD", Name: "US Dollar", Symbol: "$", DecimalPlaces: 2},
	{Code: "EUR", Name: "Euro", Symbol: "€", DecimalPlaces: 2},
	{Code: "GBP", Name: "British Pound", Symbol: "£", DecimalPlaces: 2},
	{Code: "JPY", Name: "Japanese Yen", Symbol: "¥", DecimalPlaces: 0},
	{Code: "CHF", Name: "Swiss Franc", Symbol: "Fr", DecimalPlaces: 2},
	{Code: "CAD", Name: "Canadian Dollar", Symbol: "C$", DecimalPlaces: 2},
	{Code: "AUD", Name: "Australian Dollar", Symbol: "A$", DecimalPlaces: 2},
	{Code: "CNY", Name: "Chinese Yuan", Symbol: "¥", DecimalPlaces: 2},
	{Code: "INR", Name: "Indian Rupee", Symbol: "₹", DecimalPlaces: 2},
	{Code: "RUB", Name: "Russian Ruble", Symbol: "₽", DecimalPlaces: 2},
	{Code: "BRL", Name: "Brazilian Real", Symbol: "R$", DecimalPlaces: 2},
	{Code: "KRW", Name: "South Korean Won", Symbol: "₩", DecimalPlaces: 0},
	{Code: "MXN", Name: "Mexican Peso", Symbol: "Mex$", DecimalPlaces: 2},
	{Code: "SGD", Name: "Singapore Dollar", Symbol: "S$", DecimalPlaces: 2},
	{Code: "HKD", Name: "Hong Kong Dollar", Symbol: "HK$", DecimalPlaces: 2},
	{Code: "NOK", Name: "Norwegian Krone", Symbol: "kr", DecimalPlaces: 2},
	{Code: "SEK", Name: "Swedish Krona", Symbol: "kr", DecimalPlaces: 2},
	{Code: "DKK", Name: "Danish Krone", Symbol: "kr", DecimalPlaces: 2},
	{Code: "NZD", Name: "New Zealand Dollar", Symbol: "NZ$", DecimalPlaces: 2},
	{Code: "ZAR", Name: "South African Rand", Symbol: "R", DecimalPlaces: 2},
	{Code: "TRY", Name: "Turkish Lira", Symbol: "₺", DecimalPlaces: 2},
	{Code: "PLN", Name: "Polish Zloty", Symbol: "zł", DecimalPlaces: 2},
	{Code: "THB", Name: "Thai Baht", Symbol: "฿", DecimalPlaces: 2},
	{Code: "ILS", Name: "Israeli Shekel", Symbol: "₪", DecimalPlaces: 2},
	{Code: "AED", Name: "UAE Dirham", Symbol: "د.إ", DecimalPlaces: 2},
}

