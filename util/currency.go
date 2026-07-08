package util

import "golang.org/x/text/currency"

const (
	USD = "USD"
	EUR = "EUR"
	GBP = "GBP"
	THB = "THB"
)

func IsValidCurrency(currencyCode string) bool {
	_, err := currency.ParseISO(currencyCode)
	return err == nil
}

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, GBP, THB:
		return true
	default:
		return false
	}
}
