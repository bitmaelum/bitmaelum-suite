package apikey

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strings"
	"testing"
	"time"
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

func TestNewAdminKey(t *testing.T) {
	key := NewAdminKey(4 * time.Hour)
	assert.True(t, strings.HasPrefix(key.ID, "BMK-"))
	assert.True(t, key.Admin)

	assert.True(t, key.HasPermission("anything"))
	assert.True(t, key.HasPermission("doesnt-matter"))
}

func TestNewKey(t *testing.T) {
	key := NewKey([]string{"foo"}, 4*time.Hour)
	assert.True(t, strings.HasPrefix(key.ID, "BMK-"))
	assert.False(t, key.Admin)
	assert.True(t, key.HasPermission("foo"))
	assert.False(t, key.HasPermission("bar"))
}
