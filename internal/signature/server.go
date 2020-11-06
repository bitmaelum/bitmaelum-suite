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

package signature

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
)

// SignHeader will add a server signature to a message header. This can be used to proof the origin of the message
func SignServerHeader(header *message.Header) error {
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

// VerifyHeader will verify a server signature from a message header. This can be used to proof the origin of the message
func VerifyServerHeader(header message.Header) bool {
	// Fetch public key from routing
	rs := container.GetResolveService()
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
