package util

import (
	"math/rand"
)

// List of owners names for testing
var nameList = []string{"John", "Jane", "Joe", "Jill", "Jack"}

// List of currency code for testing
var currencyList = []string{"USD", "EUR", "CAD", "CNY", "JPY"}

// By default, the random number generate number between [1, 1000]
func GetRandomInt() int64 {
	return rand.Int63n(1000) + 1
}

func GetRandomIntWithRange(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}

func GetRandomStringWithLength(length int64) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func GetRandomOwnerName() string {
	return nameList[rand.Intn(len(nameList))]
}

func GetRandomCurrency() string {
	return currencyList[rand.Intn(len(currencyList))]
}
