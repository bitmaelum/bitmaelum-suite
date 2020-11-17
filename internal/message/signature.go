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
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// SignServerHeader will add a server signature to a message header. This can be used to proof the origin of the message
func SignServerHeader(header *Header) error {
	// Already signed? Then skip
	if len(header.Signatures.Server) > 0 {
		return nil
	}

	data, err := json.Marshal(header)
	if err != nil {
		return err
	}

	h := sha256.Sum256(data)
	sig, err := bmcrypto.Sign(config.Routing.PrivateKey, h[:])
	if err != nil {
		return err
	}

	header.Signatures.Server = base64.StdEncoding.EncodeToString(sig)
	return nil
}

// VerifyServerHeader will verify a server signature from a message header. This can be used to proof the origin of the message
func VerifyServerHeader(header Header) bool {
	// Fetch public key from routing
	rs := container.Instance.GetResolveService()
	addr, err := rs.ResolveAddress(header.From.Addr)
	if err != nil {
		return false
	}

	// No header at all
	if len(header.Signatures.Server) == 0 {
		return false
	}

	// Store signature
	targetSignature, err := base64.StdEncoding.DecodeString(header.Signatures.Server)
	if err != nil {
		return false
	}
	header.Signatures.Server = ""

	// Generate hash
	data, err := json.Marshal(&header)
	if err != nil {
		return false
	}
	h := sha256.Sum256(data)

	// Verify signature
	ok, err := bmcrypto.Verify(addr.RoutingInfo.PublicKey, h[:], []byte(targetSignature))
	if err != nil {
		return false
	}

	return ok
}

// SignClientHeader will add a client signature to a message header. This can be used to proof the origin of the message
func SignClientHeader(header *Header, privKey bmcrypto.PrivKey) error {
	// Already signed? Then skip
	if len(header.Signatures.Client) > 0 {
		fmt.Println("already signed")
		return nil
	}

	data, err := json.Marshal(header)
	if err != nil {
		return err
	}

	h := sha256.Sum256(data)
	sig, err := bmcrypto.Sign(privKey, h[:])
	if err != nil {
		return err
	}

	header.Signatures.Client = base64.StdEncoding.EncodeToString(sig)
	return nil
}

// VerifyClientHeader will verify a client signature from a message header. This can be used to proof the origin of the message
func VerifyClientHeader(header Header) bool {
	// Fetch public key from routing
	rs := container.Instance.GetResolveService()
	addr, err := rs.ResolveAddress(header.From.Addr)
	if err != nil {
		return false
	}

	// No header at all
	if len(header.Signatures.Client) == 0 {
		return false
	}

	// Store signature
	targetSignature, err := base64.StdEncoding.DecodeString(header.Signatures.Client)
	if err != nil {
		return false
	}
	header.Signatures.Server = ""
	header.Signatures.Client = ""

	// Generate hash
	data, err := json.Marshal(&header)
	if err != nil {
		return false
	}
	h := sha256.Sum256(data)

	// If we have sent an authorized key key, we need to validate this first
	pubKey := addr.PublicKey
	if header.From.SignedBy == SignedByTypeAuthorized {
		pubKey = *header.AuthorizedBy.PublicKey

		// Verify our authorized key
		msg := hash.New(header.AuthorizedBy.PublicKey.String())
		sig, err := base64.StdEncoding.DecodeString(header.AuthorizedBy.Signature)
		if err != nil {
			return false
		}
		ok, err := bmcrypto.Verify(addr.PublicKey, msg.Byte(), sig)
		if err != nil || !ok {
			// Cannot validate the authorized key
			return false
		}
	}

	// Verify signature
	ok, err := bmcrypto.Verify(pubKey, h[:], []byte(targetSignature))
	if err != nil {
		return false
	}

	return ok
}
