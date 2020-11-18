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
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/stretchr/testify/assert"
)

func TestGenerateED25519KeypairFromMnemonic(t *testing.T) {
	kt, err := bmcrypto.FindKeyType("ed25519")
	assert.NoError(t, err)
	s, priv1, pub1, err := GenerateKeypairWithMnemonic(kt)
	assert.NoError(t, err)

	priv2, pub2, err := GenerateKeypairFromMnemonic(s)
	assert.NoError(t, err)

	assert.Equal(t, priv1, priv2)
	assert.Equal(t, pub1, pub2)
}

func TestGenerateRSAKeypairFromMnemonic(t *testing.T) {
	kt, err := bmcrypto.FindKeyType("rsa")
	assert.NoError(t, err)
	s, priv1, pub1, err := GenerateKeypairWithMnemonic(kt)
	assert.NoError(t, err)

	priv2, pub2, err := GenerateKeypairFromMnemonic(s)
	assert.NoError(t, err)

	assert.Equal(t, priv1, priv2)
	assert.Equal(t, pub1, pub2)
}
