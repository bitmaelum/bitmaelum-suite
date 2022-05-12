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

package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHash(t *testing.T) {
	var h Hash

	h = New("joshua!")
	assert.Equal(t, "2f92571b5567b4557b94ac5701fc48e552ba9970d6dac89f7c2ebce92f1cd836", h.String())

	h = New("somethingelse")
	assert.Equal(t, "102d0177f8ce6b0f5ada79780043ec90440529a53c98bd6043419e64c1d4274f", h.String())
}

func TestHashFromString(t *testing.T) {
	var (
		h   *Hash
		err error
	)

	h, err = NewFromHash("12345678")
	assert.Error(t, err)
	assert.Nil(t, h)

	h, err = NewFromHash("102d0177f8ce6b0f5ada79780043ec90440529a53c98bd6043419e64c1d4274f")
	assert.NoError(t, err)
	assert.Equal(t, "102d0177f8ce6b0f5ada79780043ec90440529a53c98bd6043419e64c1d4274f", h.String())
}

func TestVerify(t *testing.T) {
	// joshua@bitmaelum!
	h, _ := NewFromHash("6b024a4e51c0c4a30c3750115c66be776253880bb4af0f313e3bf2236e808840")
	assert.True(t, h.Verify(
		"fc52fabe94c0e037d2df4498e87481a6438960c9f73d517584a7a5c564535ac4",
		"49aa67181f4a3176f9b65605390bb81126e8ff1f6d03b1bd220c53e7a6b36d3e",
	))
	assert.False(t, h.Verify(
		"fc52fabe94c0e037d2df4498e87481a6438960c9f73d517584a7a5c564535ac4",
		"0000000000000006f9b65605390bb81126e8ff1f6d03b1bd220c53e7a6b36d3e",
	))

	// joshua!
	h, _ = NewFromHash("66c94b6643ada5661b2d940eb87502b5af0f47f40fd45ce0fa125502dfa9c1ee")
	assert.True(t, h.Verify(
		"fc52fabe94c0e037d2df4498e87481a6438960c9f73d517584a7a5c564535ac4",
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	))
	assert.False(t, h.Verify(
		"fc52fabe94c0e037d2df4498e87481a6438960c9f73d517584a7a5c564535ac4",
		"00000000000c1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	))
}

func TestEmpty(t *testing.T) {
	h := New("")
	assert.True(t, h.IsEmpty())

	h = New("foo")
	assert.False(t, h.IsEmpty())
}
