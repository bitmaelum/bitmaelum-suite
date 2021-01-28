// Copyright (c) 2021 BitMaelum Authors
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

package message

import (
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// ChecksumList is a list of key/value pairs of checksums. ie: ["sha1"] = "123456abcde"
type ChecksumList map[string]string

// SignedByType is a type that tells us how a message is signed
type SignedByType string

const (
	// SignedByTypeOrigin signed by origin address / private key
	SignedByTypeOrigin SignedByType = "origin"
	// SignedByTypeAuthorized signed by an authorized private key (info stored in authorizedPublicKey)
	SignedByTypeAuthorized SignedByType = "authorized"
	// SignedByTypeServer signed by the server (postmaster)
	SignedByTypeServer SignedByType = "server"
)

// AuthorizedByType holds info about the authorized sender in case the message is send and signed by an authorized sender instead of the origin sender
type AuthorizedByType struct {
	PublicKey *bmcrypto.PubKey `json:"public_key"` // Public key of the authorized sender
	Signature string           `json:"signature"`  // Signature signed by the origin address
}

// Header represents a message header
type Header struct {
	// Information on the sender of the message
	From struct {
		Addr     hash.Hash    `json:"address"`   // Address hash of the sender
		SignedBy SignedByType `json:"signed_by"` // Who has signed this message (the originator, or an authorized sender?)
	} `json:"from"`

	// Information on the recipient of the message
	To struct {
		Addr        hash.Hash `json:"address"`               // Address hash of the recipient
		Fingerprint string    `json:"fingerprint,omitempty"` // The fingerprint used for encrypting to this user
	} `json:"to"`

	// Information about the catalog of this message
	Catalog struct {
		Size          uint64       `json:"size"`           // Size of the catalog file
		Checksum      ChecksumList `json:"checksum"`       // Checksum of the catalog file
		TransactionID string       `json:"txid,omitempty"` // Transaction ID (if used) for encryption
		EncryptedKey  []byte       `json:"encrypted_key"`  // The actual encrypted key, only to be decrypted by the private key of the recipient
	} `json:"catalog"`

	AuthorizedBy *AuthorizedByType `json:"authorized_by,omitempty"` // Using a pointer type since this section can be completely omitted

	// Signatures on the message header
	Signatures struct {
		Server string `json:"server"` // Signature of the server AND client section, as signed by the private key of the server. Filled in when sending the message
		Client string `json:"client"` // Signature of the client section, as signed by the private key of the client. Filled in by the client sending the message
	} `json:"signatures"`
}
