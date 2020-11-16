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

package bmcrypto

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"errors"
)

var (
	curveFunc = elliptic.P384     // Curve used for ECDSA
	RsaBits   = []int{2048, 4096} // RSA key sizes
)

// GenerateKeyPair generates a private/public keypair based on the given type
func GenerateKeyPair(kt KeyType) (*PrivKey, *PubKey, error) {
	switch kt {
	case KeyTypeRSA:
		return generateKeyPairRSA(0)
	case KeyTypeRSAV1:
		return generateKeyPairRSA(1)
	case KeyTypeECDSA:
		return generateKeyPairECDSA()
	case KeyTypeED25519:
		return generateKeyPairED25519()
	}

	return nil, nil, errors.New("incorrect key type specified")
}

func generateKeyPairRSA(version int) (*PrivKey, *PubKey, error) {
	privRSAKey, err := rsa.GenerateKey(randReader, RsaBits[version])
	if err != nil {
		return nil, nil, err
	}

	privKey, err := NewPrivKeyFromInterface(privRSAKey)
	if err != nil {
		return nil, nil, err
	}
	pubKey, err := NewPubKeyFromInterface(privKey.K.(*rsa.PrivateKey).Public())
	if err != nil {
		return nil, nil, err
	}

	// Set the correct version of the key, since rsa can have multiple versions
	switch version {
	case 0:
		privKey.Type = KeyTypeRSA
		pubKey.Type = KeyTypeRSA
	case 1:
		privKey.Type = KeyTypeRSAV1
		pubKey.Type = KeyTypeRSAV1
	}

	return privKey, pubKey, nil
}

func generateKeyPairECDSA() (*PrivKey, *PubKey, error) {
	privECDSAKey, err := ecdsa.GenerateKey(curveFunc(), randReader)
	if err != nil {
		return nil, nil, err
	}

	privKey, err := NewPrivKeyFromInterface(privECDSAKey)
	if err != nil {
		return nil, nil, err
	}
	pubKey, err := NewPubKeyFromInterface(privKey.K.(*ecdsa.PrivateKey).Public())
	if err != nil {
		return nil, nil, err
	}

	return privKey, pubKey, nil
}

func generateKeyPairED25519() (*PrivKey, *PubKey, error) {
	pubKey, privKey, err := ed25519.GenerateKey(randReader)
	if err != nil {
		return nil, nil, err
	}

	priv, err := NewPrivKeyFromInterface(privKey)
	if err != nil {
		return nil, nil, err
	}
	pub, err := NewPubKeyFromInterface(pubKey)
	if err != nil {
		return nil, nil, err
	}
	return priv, pub, nil
}
