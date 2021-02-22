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

package bmcrypto

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomSeedToMnemonic(t *testing.T) {
	seed, _ := hex.DecodeString("cf8952d5ce19cf8d606ebf6222520a50df7bd81d33be4453")
	mnemonic, err := RandomSeedToMnemonic(seed)
	assert.NoError(t, err)
	assert.Equal(t, "sort enhance rely order ostrich shop like subject giraffe barely little payment waste ugly inquiry jelly dust occur", mnemonic)

	tmp, err := MnemonicToRandomSeed(mnemonic)
	assert.NoError(t, err)
	assert.Equal(t, "cf8952d5ce19cf8d606ebf6222520a50df7bd81d33be4453", hex.EncodeToString(tmp))

	tmp, err = MnemonicToRandomSeed("foo")
	assert.Error(t, err)
	assert.Nil(t, tmp)

	tmp, err = MnemonicToRandomSeed("ed25519 foo bar baz")
	assert.Error(t, err)
	assert.Nil(t, tmp)

	tmp, err = MnemonicToRandomSeed("ed25519 ship")
	assert.Error(t, err)
	assert.Nil(t, tmp)
}

func TestGetMnemonic(t *testing.T) {
	kp := KeyPair{
		Generator:   "FOOBAR",
		FingerPrint: "",
		PrivKey:     PrivKey{},
		PubKey:      PubKey{},
	}

	s := GetMnemonic(&kp)
	assert.Empty(t, s)
}

func TestGenerateKeypairFromMnemonic(t *testing.T) {
	kp, err := GenerateKeypairFromMnemonic("foobar")
	assert.Error(t, err)
	assert.Nil(t, kp)

	kp, err = GenerateKeypairFromMnemonic("ed25519 foobar")
	assert.Error(t, err)
	assert.Nil(t, kp)

}

func TestGenerateKeypair(t *testing.T) {
	kt, err := FindKeyType("ed25519")
	assert.NoError(t, err)

	kp, err := GenerateKeypairWithRandomSeed(kt)
	assert.NoError(t, err)

	mnemonic := GetMnemonic(kp)
	assert.Contains(t, mnemonic, "ed25519")

	kp2, err := GenerateKeypairFromMnemonic(mnemonic)
	assert.NoError(t, err)
	assert.Equal(t, kp.FingerPrint, kp2.FingerPrint)
}
