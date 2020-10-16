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
	"crypto/ed25519"
	"crypto/sha256"
	"io"

	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/hkdf"
)

// GenerateKeypairFromMnemonic generates a keypair based on the given mnemonic
func GenerateKeypairFromMnemonic(mnemonic string) (*PrivKey, *PubKey, error) {
	e, err := bip39.MnemonicToByteArray(mnemonic, true)
	if err != nil {
		return nil, nil, err
	}

	return genKey(e)
}

// GenerateKeypairWithMnemonic generates a mnemonic, and a keypair that can be generated through the same mnemonic again.
func GenerateKeypairWithMnemonic() (string, *PrivKey, *PubKey, error) {
	// Generate large enough random string
	e, err := bip39.NewEntropy(192)
	if err != nil {
		return "", nil, nil, err
	}

	// Generate Mnemonic words
	mnemonic, err := bip39.NewMnemonic(e)
	if err != nil {
		return "", nil, nil, err
	}

	privKey, pubKey, err := genKey(e)
	if err != nil {
		return "", nil, nil, err
	}

	return mnemonic, privKey, pubKey, nil
}

func genKey(e []byte) (*PrivKey, *PubKey, error) {
	// Stretch 192 bits to 256 bits
	rd := hkdf.New(sha256.New, e, []byte{}, []byte{})
	expbuf := make([]byte, 32)
	_, err := io.ReadFull(rd, expbuf)
	if err != nil {
		return nil, nil, err
	}

	// Generate keypair
	r := ed25519.NewKeyFromSeed(expbuf[:32])
	privKey, err := NewPrivKeyFromInterface(r)
	if err != nil {
		return nil, nil, err
	}
	pubKey, err := NewPubKeyFromInterface(r.Public())
	if err != nil {
		return nil, nil, err
	}

	return privKey, pubKey, nil
}
