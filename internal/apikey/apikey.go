// Copyright (c) 2020 BitMaelum Authors
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

package apikey

import (
	"math/rand"
	"strings"
	"time"
)

// KeyType represents a key with a validity and permissions. When admin is true, permission checks are always true
type KeyType struct {
	ID          string    `json:"key"`
	ValidUntil  time.Time `json:"valid_until"`
	Permissions []string  `json:"permissions"`
	Admin       bool      `json:"admin"`
}

// NewAdminKey creates a new admin key
func NewAdminKey(valid time.Duration) KeyType {
	var until = time.Time{}
	if valid > 0 {
		until = time.Now().Add(valid)
	}

	return KeyType{
		ID:          GenerateKey("BMK-", 32),
		ValidUntil:  until,
		Permissions: nil,
		Admin:       true,
	}
}

// NewKey creates a new key with the given permissions and duration
func NewKey(perms []string, valid time.Duration) KeyType {
	var until = time.Time{}
	if valid > 0 {
		until = time.Now().Add(valid)
	}

	return KeyType{
		ID:          GenerateKey("BMK-", 32),
		ValidUntil:  until,
		Permissions: perms,
		Admin:       false,
	}
}

// HasPermission returns true when this key contains the given permission, or when the key is an admin key (granted all permissions)
func (key *KeyType) HasPermission(perm string) bool {
	if key.Admin {
		return true
	}

	perm = strings.ToLower(perm)

	for _, p := range key.Permissions {
		if strings.ToLower(p) == perm {
			return true
		}
	}

	return false
}

// Repository is a repository to fetch and store API keys
type Repository interface {
	Fetch(ID string) (*KeyType, error)
	Store(key KeyType) error
	Remove(ID string)
}

// GenerateKey generates a random key based on a given string length
func GenerateKey(prefix string, n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}

	return prefix + string(b)
}
