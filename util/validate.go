package util

import (
	"regexp"
	"slices"
)

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

var supportedCurrencies = []string{USD, EUR, CAD}

func IsSupportedCurrency(currency string) bool {
	return slices.Contains(supportedCurrencies, currency)
}

func IsValidEmail(email string) bool {
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailPattern)
	return re.MatchString(email)
}