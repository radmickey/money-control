package converters

// ConvertAmount converts amount from one currency to another using provided rates
// Rates are relative to base currency (e.g., USD)
// Example: if base is USD, EUR rate = 0.85 means 1 USD = 0.85 EUR
func ConvertAmount(amount float64, from, to string, rates map[string]float64) float64 {
	if from == to {
		return amount
	}

	fromRate := rates[from]
	toRate := rates[to]

	if fromRate == 0 {
		fromRate = 1
	}
	if toRate == 0 {
		toRate = 1
	}

	// Convert: amount in 'from' -> base currency -> 'to' currency
	inBase := amount / fromRate
	return inBase * toRate
}

// BuildRatesMap creates rates map with base currency = 1.0
func BuildRatesMap(baseCurrency string, rates map[string]float64) map[string]float64 {
	result := make(map[string]float64)
	result[baseCurrency] = 1.0
	for k, v := range rates {
		result[k] = v
	}
	return result
}

