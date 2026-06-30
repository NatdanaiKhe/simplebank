package util

import (
	"math/rand/v2"
)

func RandomString(n int) string {

	b := make([]byte, n)
	for i := range b {
		b[i] = byte(rand.IntN(26) + 'a')
	}
	return string(b)
}

func RandomInt(min, max int) int64 {
	return int64(rand.IntN(max-min+1) + min)
}

func RandomFloat(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	return currencies[rand.IntN(len(currencies))]
}
