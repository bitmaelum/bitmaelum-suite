package apikey

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/stretchr/testify/assert"
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

	assert.True(t, key.HasPermission("anything", nil))
	assert.True(t, key.HasPermission("doesnt-matter", nil))

	h := hash.New("foobar")
	assert.True(t, key.HasPermission("anything", &h))
	assert.True(t, key.HasPermission("doesnt-matter", &h))
}

func TestNewKey(t *testing.T) {
	key := NewKey([]string{"foo"}, 4*time.Hour, nil)
	assert.True(t, strings.HasPrefix(key.ID, "BMK-"))
	assert.False(t, key.Admin)
	assert.True(t, key.HasPermission("foo", nil))
	assert.False(t, key.HasPermission("bar", nil))


	h := hash.New("foobar")
	key = NewKey([]string{"foo"}, 4*time.Hour, &h)
	assert.True(t, key.HasPermission("foo", &h))
	assert.False(t, key.HasPermission("bar", &h))

	h2 := hash.New("anotherhash")
	assert.False(t, key.HasPermission("foo", &h2))
	assert.False(t, key.HasPermission("bar", &h2))
}

