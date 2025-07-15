package util

var supportedCurrencies = []string{"USD", "EUR", "RUB", "KZT"}

func IsValidCurrency(currency string) bool {
	for _, c := range supportedCurrencies {
		if c == currency {
			return true
		}
	}
	return false
}
