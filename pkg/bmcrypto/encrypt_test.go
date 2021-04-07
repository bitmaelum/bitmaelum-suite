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
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	// RSA Encryption
	data, _ := ioutil.ReadFile("../../testdata/pubkey.rsa")
	pubKey, _ := PublicKeyFromString(string(data))

	data, _ = ioutil.ReadFile("../../testdata/privkey.rsa")
	privKey, _ := PrivateKeyFromString(string(data))

	cipherText, settings, err := Encrypt(*pubKey, []byte("foobar"))
	assert.NoError(t, err)
	assert.Equal(t, RsaOAEP, settings.Type)
	assert.Empty(t, settings.TransactionID)
	assert.NotEqual(t, []byte("foobar"), cipherText)

	plaintext, err := Decrypt(*privKey, settings, cipherText)
	assert.NoError(t, err)
	assert.Equal(t, []byte("foobar"), plaintext)

	// Cant decode with different type
	settings.Type = Rsav15
	plaintext, err = Decrypt(*privKey, settings, cipherText)
	assert.Error(t, err)
	assert.Empty(t, plaintext)

	// ED25519 Dual Key-Exchange + Encryption
	priv25519Key, _ := PrivateKeyFromString("ed25519 MC4CAQAwBQYDK2VwBCIEIBJsN8lECIdeMHEOZhrdDNEZl5BuULetZsbbdsZBjZ8a")
	pub25519Key, _ := PublicKeyFromString("ed25519 MCowBQYDK2VwAyEAblFzZuzz1vItSqdHbr/3DZMYvdoy17ALrjq3BM7kyKE=")
	cipher, settings, err := pub25519Key.Type.Encrypt(*pub25519Key, []byte("foobar"))
	assert.Nil(t, err)
	assert.Equal(t, Ed25519AES, settings.Type)
	assert.NotEqual(t, []byte("foobar"), cipher)

	plaintext, err = pub25519Key.Type.Decrypt(*priv25519Key, settings, cipher)
	assert.Nil(t, err)
	assert.Equal(t, []byte("foobar"), plaintext)
}
