package api

import (
	m "github.com/Petatron/bank-simulator-backend/model"
	"github.com/go-playground/validator/v10"
)

// validCurrency is a validator.Func that checks if the currency is valid
var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(m.CurrencyType); ok {
		return currency.IsValid()
	}
	return false
}
