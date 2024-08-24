package utils

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// Generates a random integer.
func RandomInt(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}

// Generates a random string of specified length.
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// Generates a random owner name.
func RandomOwner() string {
	return RandomString(6)
}

// Generates a random amount of money.
func RandomMoneyAmount() int64 {
	return RandomInt(0, 10000)
}

// Generates a random amount of money (positive or negative)
func RandomMoneyAmountForEntries() int64 {
	negative := rand.Intn(2) == 1
	if negative {
		return RandomMoneyAmount() * -1
	}
	return RandomMoneyAmount()
}

// Generates a random currency code.
func RandomCurrency() string {
	currencies := []string{EUR, USD, CAD, AUD}

	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// Generates a random email address.
func RandomEmailAddress() string {
	return fmt.Sprintf("%s@example.com", RandomString(6))
}
