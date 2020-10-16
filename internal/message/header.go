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

package message

import (
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"

	pow "github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

// ChecksumList is a list of key/value pairs of checksums. ie: ["sha1"] = "123456abcde"
type ChecksumList map[string]string

// Header represents a message header
type Header struct {
	From struct {
		Addr        hash.Hash        `json:"address"`
		PublicKey   *bmcrypto.PubKey `json:"public_key"`
		ProofOfWork *pow.ProofOfWork `json:"proof_of_work"`
	} `json:"from"`
	To struct {
		Addr hash.Hash `json:"address"`
	} `json:"to"`
	Catalog struct {
		Size          uint64       `json:"size"`
		Checksum      ChecksumList `json:"checksum"`
		Crypto        string       `json:"crypto"`
		TransactionID string       `json:"txid"`
		EncryptedKey  []byte       `json:"encrypted_key"`
	} `json:"catalog"`

	// Signature of the from, to and catalog structures, as signed by the private key of the server.
	ServerSignature string `json:"sender_signature,omitempty"`

	// Signature of the from, to and catalog structures, as signed by the private key of the client.
	ClientSignature string `json:"client_signature,omitempty"`
}

// Checksum holds a checksum which consists of the checksum hash value, and the given type of the checksum
type Checksum struct {
	Hash  string `json:"hash"`
	Value string `json:"value"`
}
