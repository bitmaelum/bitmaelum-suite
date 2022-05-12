// Copyright (c) 2022 BitMaelum Authors
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
	"crypto/elliptic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyEcdsa_NewEcdsaKey(t *testing.T) {
	kt := NewEcdsaKey(elliptic.P224())

	assert.False(t, kt.CanEncrypt())
	assert.True(t, kt.CanKeyExchange())
	assert.False(t, kt.CanDualKeyExchange())
	assert.Equal(t, "ecdsa", kt.String())

	assert.Equal(t, "ES384", kt.JWTSignMethod().Alg())

	privKey, pubKey, err := GenerateKeyPair(NewEd25519Key())
	assert.NoError(t, err)

	msg := []byte("secretmessage")

	b, settings, err := kt.Encrypt(*pubKey, msg)
	assert.NoError(t, err)
	assert.NotNil(t, settings.TransactionID)
	assert.NotEqual(t, msg, b)
	assert.Equal(t, "ecdsa+aes", string(settings.Type))

	b, err = kt.Decrypt(*privKey, settings, b)
	assert.NoError(t, err)
	assert.Equal(t, b, msg)
}

func TestKeyEcdsa_DualKeyExchange(t *testing.T) {
	kt := NewEcdsaKey(elliptic.P224())

	_, pubKey, err := kt.GenerateKeyPair(randReader)
	assert.NotNil(t, pubKey)
	assert.NoError(t, err)

	b, txid, err := kt.DualKeyExchange(*pubKey)
	assert.Error(t, err, errCannotuseForDualKeyExchange)
	assert.Nil(t, b)
	assert.Nil(t, txid)
}
