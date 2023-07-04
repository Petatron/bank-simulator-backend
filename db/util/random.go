package util

import (
	"math/rand"
)

// List of owners names for testing
var nameList = []string{"John", "Jane", "Joe", "Jill", "Jack"}

// List of currency code for testing
var currencyList = []string{"USD", "EUR", "CAD", "CNY", "JPY"}

// By default, the random number generate number between [1, 1000]
func getRandomInt() int {
	return 1 + rand.Intn(1000)
}

func getRandomIntWithRange(min, max int) int {
	return min + rand.Intn(max-min)
}

func getRandomStringWithLength(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func getRandomOwnerName() string {
	return nameList[rand.Intn(len(nameList))]
}

func getRandomCurrency() string {
	return currencyList[rand.Intn(len(currencyList))]
}
