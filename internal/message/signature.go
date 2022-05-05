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
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/sirupsen/logrus"
)

// SignServerHeader will add a server signature to a message header. This can be used to proof the origin of the message
func SignServerHeader(header *Header) error {
	h, err := generateServerHash(header)
	if err != nil {
		return err
	}

	sig, err := bmcrypto.Sign(config.Routing.KeyPair.PrivKey, h)
	if err != nil {
		return err
	}

	header.Signatures.Server = base64.StdEncoding.EncodeToString(sig)

	return nil
}

// VerifyServerHeader will verify a server signature from a message header. This can be used to proof the origin of the message
func VerifyServerHeader(header Header) bool {
	// If sent from the server there is no server signature
	if header.From.SignedBy == SignedByTypeServer {
		logrus.Trace("verifysig: signed by server. No need for checking server sig (already done on client signature)")
		return true
	}

	// Fetch public key from routing
	rs := container.Instance.GetResolveService()
	addr, err := rs.ResolveAddress(header.From.Addr)
	if err != nil {
		logrus.Trace("verifysig: cant find header from addr")
		return false
	}

	// No header at all
	if len(header.Signatures.Server) == 0 {
		logrus.Trace("verifysig: no header signature found")
		return false
	}

	// Generate server hash
	h, err := generateServerHash(&header)
	if err != nil {
		logrus.Trace("verifysig: error generating server hash")
		return false
	}

	// Decode signature
	targetSignature, err := base64.StdEncoding.DecodeString(header.Signatures.Server)
	if err != nil {
		logrus.Trace("verifysig: error decoding target signature")
		return false
	}

	// Verify signature
	ok, err := bmcrypto.Verify(addr.RoutingInfo.PublicKey, h, targetSignature)
	if err != nil {
		logrus.Trace("verifysig: error verifying with bmcrypto")
		return false
	}

	return ok
}

// SignClientHeader will add a client signature to a message header. This can be used to proof the origin of the message
func SignClientHeader(header *Header, privKey bmcrypto.PrivKey) error {

	// Generate client hash
	h, err := generateClientHash(header)
	if err != nil {
		return nil
	}

	// Sign
	sig, err := bmcrypto.Sign(privKey, h)
	if err != nil {
		return err
	}

	header.Signatures.Client = base64.StdEncoding.EncodeToString(sig)
	return nil
}

// VerifyClientHeader will verify a client signature from a message header. This can be used to proof the origin of the message
func VerifyClientHeader(header Header) bool { // No header at all
	if len(header.Signatures.Client) == 0 {
		return false
	}

	// Store signature
	targetSignature, err := base64.StdEncoding.DecodeString(header.Signatures.Client)
	if err != nil {
		return false
	}

	// Generate client hash
	h, err := generateClientHash(&header)
	if err != nil {
		return false
	}

	signedByPublicKey, err := getSignerPublicKey(header)
	if err != nil {
		return false
	}

	// Verify signature
	ok, err := bmcrypto.Verify(*signedByPublicKey, h[:], []byte(targetSignature))
	if err != nil {
		return false
	}

	return ok
}

func getSignerPublicKey(header Header) (*bmcrypto.PubKey, error) {
	switch header.From.SignedBy {
	case SignedByTypeServer:
		// Resolve the routing to fetch the public key since the From Addr is the routing ID
		rs := container.Instance.GetResolveService()
		routing, err := rs.ResolveRouting(header.From.Addr.String())
		if err != nil {
			return nil, err
		}
		return &routing.PublicKey, err

	case SignedByTypeAuthorized:
		// Fetch public key from routing
		rs := container.Instance.GetResolveService()
		addr, err := rs.ResolveAddress(header.From.Addr)
		if err != nil {
			return nil, err
		}

		// Check if the authorized_by is filled
		if header.AuthorizedBy == nil || header.AuthorizedBy.PublicKey == nil || header.AuthorizedBy.Signature == "" {
			return nil, errors.New("no authorization key found")
		}

		msg := hash.New(header.AuthorizedBy.PublicKey.String())
		sig, err := base64.StdEncoding.DecodeString(header.AuthorizedBy.Signature)
		if err != nil {
			return nil, err
		}

		// Test if the authorized public key is actually signed by the authorizer
		ok, err := bmcrypto.Verify(addr.PublicKey, msg.Byte(), sig)
		if err != nil || !ok {
			// Cannot validate the authorized key
			return nil, err
		}

		// The signature is correct (the key is signed by the originating authorizer). The can safely be used for verifying our client signature
		return header.AuthorizedBy.PublicKey, nil

	default:
		// Fetch public key from routing
		rs := container.Instance.GetResolveService()
		addr, err := rs.ResolveAddress(header.From.Addr)
		if err != nil {
			return nil, err
		}
		return &addr.PublicKey, nil
	}
}

func generateServerHash(header *Header) ([]byte, error) {
	h, err := generateClientHash(header)
	if err != nil {
		return nil, err
	}

	hdr := map[string]interface{}{
		"client_hash":          h,
		"server_authorized_by": header.AuthorizedBy,
	}

	data, err := json.Marshal(hdr)
	if err != nil {
		return nil, err
	}

	h2 := sha256.Sum256(data)
	return h2[:], nil
}

func generateClientHash(header *Header) ([]byte, error) {
	hdr := map[string]interface{}{
		"client_from":    header.From,
		"client_to":      header.To,
		"client_catalog": header.Catalog,
	}

	data, err := json.Marshal(hdr)
	if err != nil {
		return nil, err
	}

	h := sha256.Sum256(data)
	return h[:], nil
}
