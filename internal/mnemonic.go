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

package internal

import (
	"encoding/hex"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/tyler-smith/go-bip39"
)

// GetMnemonic will return a mnemonic representation of the given keypair
func GetMnemonic(kp *bmcrypto.KeyPair) string {
	seed, err := hex.DecodeString(kp.Generator)
	if err != nil {
		return ""
	}

	mnemonic, err := RandomSeedToMnemonic(seed)
	if err != nil {
		return ""
	}

	return kp.PubKey.Type.String() + " " + mnemonic
}

// RandomSeedToMnemonic converts a random seed to a mnemonic
func RandomSeedToMnemonic(seed []byte) (string, error) {
	return bip39.NewMnemonic(seed)
}

// MnemonicToRandomSeed converts a mnemonic to a random seed
func MnemonicToRandomSeed(mnemonic string) ([]byte, error) {
	return bip39.MnemonicToByteArray(mnemonic, true)
}

// GenerateKeypairWithRandomSeed generates a seed and generates a keypair which can be reconstructed again with the same seed
func GenerateKeypairWithRandomSeed(kt bmcrypto.KeyType) (*bmcrypto.KeyPair, error) {
	// Generate large enough random string
	seed, err := bip39.NewEntropy(192)
	if err != nil {
		return nil, err
	}

	return bmcrypto.CreateKeypair(kt, seed)
}

// GenerateKeypairFromRandomSeed generates a keypair based on the given mnemonic
func GenerateKeypairFromMnemonic(mnemonic string) (*bmcrypto.KeyPair, error) {
	words := strings.SplitN(mnemonic, " ", 2)

	kt, err := bmcrypto.FindKeyType(words[0])
	if err != nil {
		return nil, err
	}

	seed, err := MnemonicToRandomSeed(words[1])
	if err != nil {
		return nil, err
	}

	return bmcrypto.CreateKeypair(kt, seed)
}
