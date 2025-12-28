package utils

import (
	"regexp"
	"strings"
)

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	currencyRegex = regexp.MustCompile(`^[A-Z]{3}$`)
	uuidRegex     = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

// IsValidEmail checks if email is valid
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// IsValidCurrency checks if currency code is valid (3-letter uppercase)
func IsValidCurrency(currency string) bool {
	return currencyRegex.MatchString(currency)
}

// IsValidUUID checks if string is a valid UUID
func IsValidUUID(id string) bool {
	return uuidRegex.MatchString(strings.ToLower(id))
}

// SupportedCurrencies is a list of commonly supported currencies
var SupportedCurrencies = map[string]string{
	"USD": "US Dollar",
	"EUR": "Euro",
	"GBP": "British Pound",
	"JPY": "Japanese Yen",
	"CHF": "Swiss Franc",
	"CAD": "Canadian Dollar",
	"AUD": "Australian Dollar",
	"NZD": "New Zealand Dollar",
	"CNY": "Chinese Yuan",
	"HKD": "Hong Kong Dollar",
	"SGD": "Singapore Dollar",
	"INR": "Indian Rupee",
	"RUB": "Russian Ruble",
	"BRL": "Brazilian Real",
	"MXN": "Mexican Peso",
	"KRW": "South Korean Won",
	"TRY": "Turkish Lira",
	"ZAR": "South African Rand",
	"SEK": "Swedish Krona",
	"NOK": "Norwegian Krone",
	"DKK": "Danish Krone",
	"PLN": "Polish Zloty",
	"CZK": "Czech Koruna",
	"HUF": "Hungarian Forint",
	"ILS": "Israeli Shekel",
	"THB": "Thai Baht",
	"PHP": "Philippine Peso",
	"IDR": "Indonesian Rupiah",
	"MYR": "Malaysian Ringgit",
	"AED": "UAE Dirham",
	"SAR": "Saudi Riyal",
}

// IsSupportedCurrency checks if currency is in supported list
func IsSupportedCurrency(currency string) bool {
	_, ok := SupportedCurrencies[strings.ToUpper(currency)]
	return ok
}

// ValidateRequired checks if a string is not empty
func ValidateRequired(value, field string) (string, bool) {
	if strings.TrimSpace(value) == "" {
		return field + " is required", false
	}
	return "", true
}

// ValidateMinLength checks if a string meets minimum length
func ValidateMinLength(value string, min int, field string) (string, bool) {
	if len(value) < min {
		return field + " must be at least " + string(rune(min+'0')) + " characters", false
	}
	return "", true
}

// ValidateMaxLength checks if a string doesn't exceed maximum length
func ValidateMaxLength(value string, max int, field string) (string, bool) {
	if len(value) > max {
		return field + " must be at most " + string(rune(max+'0')) + " characters", false
	}
	return "", true
}

// ValidatePositive checks if a number is positive
func ValidatePositive(value float64, field string) (string, bool) {
	if value <= 0 {
		return field + " must be positive", false
	}
	return "", true
}

// ValidateNonNegative checks if a number is non-negative
func ValidateNonNegative(value float64, field string) (string, bool) {
	if value < 0 {
		return field + " cannot be negative", false
	}
	return "", true
}

