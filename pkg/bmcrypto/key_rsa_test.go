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

func TestString(t *testing.T) {
	k := NewRsaKey(2048)
	assert.Equal(t, "rsa", k.String())

	k = NewRsaKey(4096)
	assert.Equal(t, "rsa4096", k.String())

	k = NewRsaKey(2048)
	assert.Equal(t, "RS256", k.JWTSignMethod().Alg())

	data, err := ioutil.ReadFile("../../testdata/privkey.rsa")
	assert.NoError(t, err)
	privKey, err := PrivateKeyFromString(string(data))
	assert.NoError(t, err)
	data, err = ioutil.ReadFile("../../testdata/pubkey.rsa")
	assert.NoError(t, err)
	pubKey, err := PublicKeyFromString(string(data))
	assert.NoError(t, err)

	b, err := k.KeyExchange(*privKey, *pubKey)
	assert.Error(t, err)
	assert.Nil(t, b)

	b, tx, err := k.DualKeyExchange(*pubKey)
	assert.Error(t, err)
	assert.Nil(t, b)
	assert.Nil(t, tx)

}
