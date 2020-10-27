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

package authkey

import (
	"time"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// KeyType represents a public key that is authorized to send email. This public key is signed by the private key of the owner and stored in the signature.
type KeyType struct {
	Fingerprint string           `json:"fingerprint"` // Fingerprint SHA256 of the public key
	AddrHash    *hash.Hash       `json:"addr_hash"`   // Address hash of the account that is authorized on
	Expires     time.Time        `json:"expires"`     // Expiry time
	PublicKey   *bmcrypto.PubKey `json:"public_key"`  // Actual public key
	Signature   string           `json:"signature"`   // Signed public key by the owner's private key
	Description string           `json:"description"` // User description
}

// NewKey creates a new authorized key
func NewKey(addrHash hash.Hash, pubKey *bmcrypto.PubKey, signature string, expiry time.Time, desc string) KeyType {
	return KeyType{
		Fingerprint: pubKey.Fingerprint(),
		Signature:   signature,
		PublicKey:   pubKey,
		Expires:     expiry,
		AddrHash:    &addrHash,
		Description: desc,
	}
}

// Repository is a repository to fetch and store auth keys
type Repository interface {
	Fetch(fingerprint string) (*KeyType, error)
	FetchByHash(h string) ([]KeyType, error)
	Store(key KeyType) error
	Remove(key KeyType) error
}
