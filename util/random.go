package util

import (
	"math/rand/v2"
	"strings"
)

//create random string as name
//create random currency
//create random balance

const alphabet = "abcdefhijklmnopqrstuvwxyz"

// generates a random string of length n
func RandomString(n int) string {
	alphabetLen := len(alphabet)
	var sb strings.Builder

	for i := 0; i < n; i++ {
		letterIndex := rand.IntN(alphabetLen)
		letter := alphabet[letterIndex]
		sb.WriteByte(letter)
	}

	return sb.String()
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int64N(max-min+1)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(1, 1000)
}


// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	n := len(currencies)
	return currencies[rand.IntN(n)]
}

