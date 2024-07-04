package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

type CurrencyRate struct {
	Rate float64 `json:"rate"`
}

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// Random owner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

func RandomBalance(currency string) float64 {
	amounts := map[string]CurrencyRate{
		"EUR": {Rate: 1.10},   // Euros per USD
		"XAF": {Rate: 607.29}, // West African Francs per USD
		"CAD": {Rate: 1.35},   // Canadian dollar per USD
		"USD": {Rate: 1.00},   // USD per USD
	}

	return float64(RandomInt(int64(amounts[currency].Rate)*100*10, int64(amounts[currency].Rate)*100*100))
}

func RandomCurrency() string {
	currencies := []string{"USD", "XAF", "EUR", "CAD"}
	n := len(currencies)

	return currencies[rand.Intn(n)]
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(6))
}
