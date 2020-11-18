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

package internal

import (
	"bytes"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/tyler-smith/go-bip39"
)

// GenerateKeypairWithMnemonic generates a mnemonic, and a RSA keypair that can be generated through the same mnemonic again.
func GenerateKeypairWithMnemonic(kt bmcrypto.KeyType) (string, *bmcrypto.PrivKey, *bmcrypto.PubKey, error) {
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

	privKey, pubKey, err := kt.GenerateKeyPair(bytes.NewReader(e))
	if err != nil {
		return "", nil, nil, err
	}

	return kt.String() + " " + mnemonic, privKey, pubKey, nil
}

// GenerateKeypairFromMnemonic generates a keypair based on the given mnemonic
func GenerateKeypairFromMnemonic(mnemonic string) (*bmcrypto.PrivKey, *bmcrypto.PubKey, error) {
	words := strings.SplitN(mnemonic, " ", 2)

	kt, err := bmcrypto.FindKeyType(words[0])
	if err != nil {
		return nil, nil, err
	}

	e, err := bip39.MnemonicToByteArray(words[1], true)
	if err != nil {
		return nil, nil, err
	}

	return kt.GenerateKeyPair(bytes.NewReader(e))
}
