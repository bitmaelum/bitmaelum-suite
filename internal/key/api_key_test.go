// Copyright (c) 2022 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package key

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
	expiry := time.Unix(1603442983, 0)

	key := NewAPIAdminKey(expiry, "foobar")
	assert.True(t, strings.HasPrefix(key.ID, "BMK-"))
	assert.True(t, key.Admin)
	assert.Equal(t, "foobar", key.Desc)

	assert.True(t, key.HasPermission("anything", nil))
	assert.True(t, key.HasPermission("doesnt-matter", nil))

	h := hash.New("foobar")
	assert.True(t, key.HasPermission("anything", &h))
	assert.True(t, key.HasPermission("doesnt-matter", &h))
}

func TestNewKey(t *testing.T) {
	expiry := time.Unix(1603442983, 0)

	key := NewAPIKey([]string{"foo"}, expiry, "")
	assert.True(t, strings.HasPrefix(key.ID, "BMK-"))
	assert.False(t, key.Admin)
	assert.Nil(t, key.AddressHash)
	assert.True(t, key.HasPermission("foo", nil))
	assert.False(t, key.HasPermission("bar", nil))

	h := hash.New("foobar")
	key = NewAPIAccountKey(h, []string{"foo"}, expiry, "")
	assert.True(t, key.HasPermission("foo", &h))
	assert.False(t, key.HasPermission("bar", &h))

	h2 := hash.New("anotherhash")
	assert.False(t, key.HasPermission("foo", &h2))
	assert.False(t, key.HasPermission("bar", &h2))
}
