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
	"io/ioutil"
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	// RSA Encryption
	data, _ := ioutil.ReadFile("../../testdata/pubkey.rsa")
	pubKey, _ := bmcrypto.PubKeyFromString(string(data))

	data, _ = ioutil.ReadFile("../../testdata/privkey.rsa")
	privKey, _ := bmcrypto.PrivKeyFromString(string(data))

	cipher, _, c, err := Encrypt(*pubKey, []byte("foobar"))
	assert.Nil(t, err)
	assert.Equal(t, "rsa", c)
	assert.NotEqual(t, []byte("foobar"), cipher)

	plaintext, err := Decrypt(*privKey, "", cipher)
	assert.Nil(t, err)
	assert.Equal(t, []byte("foobar"), plaintext)

	// ED25519 Dual Key-Exchange + Encryption
	priv25519Key, _ := bmcrypto.PrivKeyFromString("ed25519 MC4CAQAwBQYDK2VwBCIEIBJsN8lECIdeMHEOZhrdDNEZl5BuULetZsbbdsZBjZ8a")
	pub25519Key, _ := bmcrypto.PubKeyFromString("ed25519 MCowBQYDK2VwAyEAblFzZuzz1vItSqdHbr/3DZMYvdoy17ALrjq3BM7kyKE=")
	cipher, txID, c, err := pub25519Key.Type.Encrypt(*pub25519Key, []byte("foobar"))
	assert.Nil(t, err)
	assert.Equal(t, "ed25519", c)
	assert.NotEqual(t, []byte("foobar"), cipher)

	plaintext, err = pub25519Key.Type.Decrypt(*priv25519Key, txID, cipher)
	assert.Nil(t, err)
	assert.Equal(t, []byte("foobar"), plaintext)
}

func TestErrors(t *testing.T) {
	kt, err := bmcrypto.FindKeyType("ecdsa")
	assert.NoError(t, err)
	priv, pub, err := bmcrypto.GenerateKeyPair(kt)
	assert.NoError(t, err)

	cipher, txID, c, err := Encrypt(*pub, []byte("foobar"))
	assert.Error(t, err)
	assert.Nil(t, cipher)
	assert.Equal(t, "ecdsa", c)
	assert.Equal(t, "", txID)

	_, err = Decrypt(*priv, "", []byte("foobar"))
	assert.Error(t, err)

}
