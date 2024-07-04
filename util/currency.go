package util

import (
	"fmt"
	"strconv"
)


var currencies = map[string]string{
	"EUR": "EUR",
	"XAF": "XAF",
	"CAD": "CAD",
	"USD": "USD",
}

var rates = map[string]CurrencyRate{
	"EUR": {Rate: 1.1},    // Euros per USD
	"XAF": {Rate: 607.29}, // West African Francs per USD
	"CAD": {Rate: 1.35},   // Canadian dollar per USD
	"USD": {Rate: 1},      // USD per USD
}

func IsSupportedCurrency(currency string) bool {
	_, ok := currencies[currency]
	return ok
}

func Converter(fromCurrency string, toCurrency string, amount float64) (float64, error) {

	if _, ok := rates[fromCurrency]; !ok {
		return 0, fmt.Errorf("unsupported currency: %s", fromCurrency)
	}
	if _, ok := rates[toCurrency]; !ok {
		return 0, fmt.Errorf("unsupported currency: %s", toCurrency)
	}

	if fromCurrency == toCurrency {
		return amount, nil
	}

	// Convert to USD first
	usdAmount := float64(amount) / rates[fromCurrency].Rate

	// Convert from USD to target currency
	convertedAmount := usdAmount * rates[toCurrency].Rate
	roundedAmount := fmt.Sprintf("%.2f", convertedAmount)
	parsedAmount, _ := strconv.ParseFloat(roundedAmount, 64)

	return parsedAmount, nil
}
