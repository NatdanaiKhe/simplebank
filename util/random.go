package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(r.Intn(26) + 'a')
	}
	return string(b)
}

func RandomInt(min, max int) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return int64(r.Intn(max-min+1) + min)
}

func RandomFloat(min, max float64) float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Float64()*(max-min) + min
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	return currencies[rand.Intn(len(currencies))]
}
