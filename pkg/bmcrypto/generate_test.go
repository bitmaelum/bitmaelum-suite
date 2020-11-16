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
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	message = []byte("this is the message we need to sign")
)

func TestGenerate(t *testing.T) {
	privKey, pubKey, err := GenerateKeyPair(KeyTypeRSA)
	assert.Nil(t, err)
	assert.IsType(t, (*rsa.PrivateKey)(nil), privKey.K)
	assert.IsType(t, (*rsa.PublicKey)(nil), pubKey.K)
	assert.Equal(t, pubKey.K.(*rsa.PublicKey).Size() * 8, RsaBits[0])

	// Check if we can verify with this key
	sig, err := Sign(*privKey, message)
	assert.Nil(t, err)
	b, err := Verify(*pubKey, message, sig)
	assert.Nil(t, err)
	assert.True(t, b)

	privKey, pubKey, err = GenerateKeyPair(KeyTypeECDSA)
	assert.Nil(t, err)
	assert.IsType(t, (*ecdsa.PrivateKey)(nil), privKey.K)
	assert.IsType(t, (*ecdsa.PublicKey)(nil), pubKey.K)

	// Check if we can verify with this key
	sig, err = Sign(*privKey, message)
	assert.Nil(t, err)
	b, err = Verify(*pubKey, message, sig)
	assert.Nil(t, err)
	assert.True(t, b)

	privKey, pubKey, err = GenerateKeyPair(KeyTypeED25519)
	assert.Nil(t, err)
	assert.IsType(t, (ed25519.PrivateKey)(nil), privKey.K)
	assert.IsType(t, (ed25519.PublicKey)(nil), pubKey.K)

	// Check if we can verify with this key
	sig, err = Sign(*privKey, message)
	assert.Nil(t, err)
	b, err = Verify(*pubKey, message, sig)
	assert.Nil(t, err)
	assert.True(t, b)

	// Unknown key
	_, _, err = GenerateKeyPair("foobar")
	assert.EqualError(t, err, "incorrect key type specified")
}

func TestRSAV1(t *testing.T) {
	privKey, pubKey, err := GenerateKeyPair(KeyTypeRSAV1)
	assert.Nil(t, err)
	assert.IsType(t, (*rsa.PrivateKey)(nil), privKey.K)
	assert.IsType(t, (*rsa.PublicKey)(nil), pubKey.K)
	assert.Equal(t, pubKey.K.(*rsa.PublicKey).Size() * 8, RsaBits[1])

}
