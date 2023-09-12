package util

import (
	"github.com/Petatron/bank-simulator-backend/model"
	"math/rand"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// GetRandomInt generate a random integer number
func GetRandomInt() int64 {
	return rand.Int63()
}

// GetRandomIntWithRange generate a random integer number with range [min, max)
func GetRandomIntWithRange(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}

// GetRandomMoneyAmount generate a random money amount with range [1, 1000)
func GetRandomMoneyAmount() int64 {
	return GetRandomIntWithRange(1, 1000)
}

// GetRandomStringWithLength generate a random string from alphabet with given length
func GetRandomStringWithLength(length int) string {
	tempList := make([]byte, length)
	for i := range tempList {
		tempList[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(tempList)
}

// GetRandomOwnerName generate a random owner name
func GetRandomOwnerName() string {
	// Get a random length from 3 to 5
	length := rand.Intn(6) + 6
	return GetRandomStringWithLength(length)
}

// GetRandomCurrency generate a random currency code
func GetRandomCurrency() string {
	// List of currency code for testing
	currencyList := make([]string, 0, len(model.CurrencyMap))
	for _, value := range model.CurrencyMap {
		currencyList = append(currencyList, value)
	}
	return currencyList[rand.Intn(len(currencyList))]
}

// GetRandomEmail generate a random email
func GetRandomEmail() string {
	return GetRandomStringWithLength(10) + "@" + GetRandomStringWithLength(5) + ".com"
}
