package internal

import (
	"math/rand"
	"time"
)

// GenerateKey generates a random key based on a given string length
func GenerateKey(prefix string, n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}

	return prefix + string(b)
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
