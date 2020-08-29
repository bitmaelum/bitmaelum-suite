package encrypt

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	var s string
	rand.Seed(99)

	s = GenerateKey("ABC-", 5)
	assert.Equal(t, "ABC-ZXQuE", s)

	s = GenerateKey("ABC-", 0)
	assert.Equal(t, "ABC-", s)

	s = GenerateKey("XA-", 32)
	assert.Equal(t, "XA-hdqR0iTCDT1YzUL7r81Ahy7qsvmqHl8l", s)
}
