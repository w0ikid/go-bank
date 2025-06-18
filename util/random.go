package util

import (
	_ "math/rand"
	"time"
	_ "strings"
	"github.com/brianvoe/gofakeit/v7"
)

// const alphabet = "abcdefghijklmnopqrstuvwxyz"

// func init() {
// 	rand.Seed(time.Now().UnixNano())
// }

// func RandomInt(min, max int64) int64 {
// 	return min + rand.Int63n(max - min + 1)
// }

// func RandomString(n int) string {
// 	var sb strings.Builder

// 	k := len(alphabet)

// 	for i := 0; i < n; i++ {
// 		c := alphabet[rand.Intn(k)]
// 		sb.WriteByte(c)
// 	}
// 	return sb.String()
// }

// func RandomOwner() string {
// 	return RandomString(6)
// }

// func RandomCurrency() string {
// 	currencies := []string{"USD", "EUR", "CAD", "AUD", "KZT"}
// 	n := len(currencies)
// 	return currencies[rand.Intn(n)]
// }

// func RandomBalance() int64 {
// 	return RandomInt(0, 100000)
// }


// -----------------------------------------------
// using gofakeit for more realistic data generation

func init() {
	gofakeit.Seed(time.Now().UnixNano())
}

func RandomOwner() string {
	return gofakeit.Name()
}

func RandomCurrency() string {
	return gofakeit.CurrencyShort()
}

func RandomBalance() int64 {
	return int64(gofakeit.Number(1, 1000000))
}