package model

type CurrencyType string

const (
	USD CurrencyType = "USD"
	EUR CurrencyType = "EUR"
	CAD CurrencyType = "CAD"
)

func (c CurrencyType) IsValid() bool {
	switch c {
	case USD, EUR, CAD:
		return true
	}
	return false
}
