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

package key

import (
	"time"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// AuthKeyType represents a public key that is authorized to send email. This public key is signed by the private key of the owner and stored in the signature.
type AuthKeyType struct {
	Fingerprint string           `json:"fingerprint"`  // Fingerprint / ID
	AddressHash hash.Hash        `json:"address_hash"` // Address hash
	Expires     time.Time        `json:"expires"`      // Expiry time
	PublicKey   *bmcrypto.PubKey `json:"public_key"`   // Actual public key
	Signature   string           `json:"signature"`    // Signed public key by the owner's private key
	Description string           `json:"description"`  // User description
}

// GetID returns fingerprint of the public key
func (a AuthKeyType) GetID() string {
	return a.Fingerprint
}

// GetAddressHash returns address hash
func (a AuthKeyType) GetAddressHash() *hash.Hash {
	return &a.AddressHash
}

// NewAuthKey creates a new authorized key
func NewAuthKey(addrHash hash.Hash, pubKey *bmcrypto.PubKey, signature string, expiry time.Time, desc string) AuthKeyType {
	return AuthKeyType{
		Fingerprint: pubKey.Fingerprint(),
		AddressHash: addrHash,
		Signature:   signature,
		PublicKey:   pubKey,
		Expires:     expiry,
		Description: desc,
	}
}
