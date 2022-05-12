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
	"time"

	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// APIKeyType represents a key with a validity and permissions. When admin is true, permission checks are always true
type APIKeyType struct {
	ID          string     `json:"key"`          // ID
	AddressHash *hash.Hash `json:"address_hash"` // Address hash
	Expires     time.Time  `json:"expires"`      // Time the key expires
	Permissions []string   `json:"permissions"`  // List of permissions for this key
	Admin       bool       `json:"admin"`        // Admin key or not?
	Desc        string     `json:"description"`  // Description of the key
}

// GetID returns ID
func (key *APIKeyType) GetID() string {
	return key.ID
}

// GetAddressHash returns address hash or nil
func (key *APIKeyType) GetAddressHash() *hash.Hash {
	return key.AddressHash
}

// NewAPIAdminKey creates a new admin key
func NewAPIAdminKey(expiry time.Time, desc string) APIKeyType {
	return APIKeyType{
		ID:          GenerateKey("BMK-", 32),
		AddressHash: nil,
		Expires:     expiry,
		Permissions: nil,
		Admin:       true,
		Desc:        desc,
	}
}

// NewAPIAccountKey creates a new key with the given permissions and duration for the given account
func NewAPIAccountKey(addrHash hash.Hash, perms []string, expiry time.Time, desc string) APIKeyType {
	return APIKeyType{
		ID:          GenerateKey("BMK-", 32),
		AddressHash: &addrHash,
		Expires:     expiry,
		Permissions: perms,
		Admin:       false,
		Desc:        desc,
	}
}

// NewAPIKey creates a new key with the given permissions and duration
func NewAPIKey(perms []string, expiry time.Time, desc string) APIKeyType {
	return APIKeyType{
		ID:          GenerateKey("BMK-", 32),
		AddressHash: nil,
		Expires:     expiry,
		Permissions: perms,
		Admin:       false,
		Desc:        desc,
	}
}

// HasPermission returns true when this key contains the given permission, or when the key is an admin key (granted all permissions)
func (key *APIKeyType) HasPermission(perm string, addrHash *hash.Hash) bool {
	if key.Admin {
		return true
	}

	perm = strings.ToLower(perm)

	for _, p := range key.Permissions {
		if strings.ToLower(p) == perm && (addrHash == nil || addrHash.String() == key.AddressHash.String()) {
			return true
		}
	}

	return false
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
