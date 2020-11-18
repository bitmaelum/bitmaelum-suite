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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSigningMethodEdDSAAlg(t *testing.T) {
	m := &SigningMethodEdDSA{}
	assert.Equal(t, "EdDSA", m.Alg())
}

func TestSigningMethodEdDSASign(t *testing.T) {
	m := &SigningMethodEdDSA{}

	privKey, _ := PrivKeyFromString("ed25519 MC4CAQAwBQYDK2VwBCIEILq+V/CUlMdbmoQC1odEgOEmtMBQu0UpIICxJbQM1vhd")
	pubKey, _ := PubKeyFromString("ed25519 MCowBQYDK2VwAyEARdZSwluYtMWTGI6Rvl0Bhu40RBDn6D88wyzFL1IR3DU=")

	s, err := m.Sign("foobar", privKey.K.(ed25519.PrivateKey))
	assert.NoError(t, err)

	err = m.Verify("foobar", s, pubKey.K.(ed25519.PublicKey))
	assert.NoError(t, err)

	err = m.Verify("foobarfoofofoo", s, pubKey.K.(ed25519.PublicKey))
	assert.Error(t, err)
}
