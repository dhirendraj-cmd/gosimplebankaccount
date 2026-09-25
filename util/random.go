package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabets = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.NewSource(time.Now().UnixNano())
	// source := rand.NewSource(time.Now().UnixNano())
	// localRand := rand.New(source)
}

// RandomInt generates a random integerer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// random string generator
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabets)

	for range n {
		c := alphabets[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandOwner() string {
	return RandomString(5)
}

func RandMoney() int64 {
	return RandomInt(10, 10000)
}

func RandCurrency() string {
	currencies := []string{"INR", "EUR", "USD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
