package util

import (
	"github.com/brianvoe/gofakeit/v7"
	_ "math/rand"
	_ "strings"
	"time"
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

func RandomInt(min, max int64) int64 {
	return int64(gofakeit.Number(int(min), int(max)))
}

func RandomOwner() string {
	return gofakeit.Name()
}

func RandomCurrency() string {
	return supportedCurrencies[gofakeit.Number(0, len(supportedCurrencies)-1)]
}

func RandomBalance() int64 {
	return int64(gofakeit.Number(1, 1000000))
}

func RandomEmail() string {
	return gofakeit.Email()
}

func RandomString(n int) string {
	return gofakeit.LetterN(uint(n))
}
