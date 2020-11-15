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

package hash

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
)

var (
	hashRegex = regexp.MustCompile(`^[a-f0-9]{64}$`)
)

var (
	// emptyOrgHash is a hash of an empty string. Used for checking against empty org hashes
	emptyOrgHash = New("")
)

// Hash is a hashed entity. This could be an address, localpart, organisation part or anything else. Currently only
// sha256 is supported
type Hash string

// New generates a regular hash. Assumes you know what you are hashing
func New(s string) Hash {
	/*if len(s) == 64 {
		panic("seems like we are hashing a hash...")
	}*/
	sum := sha256.Sum256([]byte(s))

	return Hash(hex.EncodeToString(sum[:]))
}

// NewFromHash generates a hash address based on the given string hash
func NewFromHash(hash string) (*Hash, error) {
	if !hashRegex.MatchString(strings.ToLower(hash)) {
		return nil, errors.New("incorrect hash address format specified")
	}

	h := Hash(hash)
	return &h, nil
}

// String casts an hash address to string
func (ha Hash) String() string {
	return string(ha)
}

// Byte casts an hash address to string
func (ha Hash) Byte() []byte {
	return []byte(ha)
}

// IsEmpty returns true when this hash is of an empty string
func (ha Hash) IsEmpty() bool {
	return ha.String() == emptyOrgHash.String()
}

// Verify will check if the hashes for local and org found matches the actual target hash
func (ha Hash) Verify(localHash, orgHash Hash) bool {
	target := New(localHash.String() + orgHash.String())

	return subtle.ConstantTimeCompare(ha.Byte(), target.Byte()) == 1
}
