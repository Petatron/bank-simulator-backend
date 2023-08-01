package model

type CurrencyType string

// Currency Types supported by the system.
const (
	USD CurrencyType = "USD"
	EUR CurrencyType = "EUR"
	CAD CurrencyType = "CAD"
)

// IsValid check if the currency type is supported.
func (c CurrencyType) IsValid() bool {
	switch c {
	case USD, EUR, CAD:
		return true
	}
	return false
}
