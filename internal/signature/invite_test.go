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

package signature

import (
	"testing"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/stretchr/testify/assert"
)

func TestSignature(t *testing.T) {
	addr := hash.New("john@acme!")

	privKey, err := bmcrypto.PrivateKeyFromString("ed25519 MC4CAQAwBQYDK2VwBCIEILq+V/CUlMdbmoQC1odEgOEmtMBQu0UpIICxJbQM1vhd")
	assert.NoError(t, err)
	pubKey, err := bmcrypto.PublicKeyFromString("ed25519 MCowBQYDK2VwAyEARdZSwluYtMWTGI6Rvl0Bhu40RBDn6D88wyzFL1IR3DU=")
	assert.NoError(t, err)

	// Assume this is the current time during tests
	timeNow = func() time.Time {
		return time.Date(2010, 01, 01, 12, 34, 56, 0, time.UTC)
	}

	it, err := NewInviteToken(addr, "12345678", time.Date(2010, 01, 04, 12, 34, 56, 0, time.UTC), *privKey)
	assert.NoError(t, err)
	assert.Equal(t, "ZTE3NmRlOTk5MDZlMjFmZWE4MWY4Yzg3YjY4MzdiNzNkM2IyMzJjNDU5MmYyMGFjMzA2NWU1ODZiNmUxN2RiNzoxMjM0NTY3ODoxMjYyNjA4NDk2OtuWWRz3/6zPsXKDQk18SEmRabib6ogVd9ml1lPYWYw4tcb940J4ZxFK77rU6rkGYf/fKG1anE1SLpUyXBNGxgM=", it.String())
	ok := it.Verify("12345678", *pubKey)
	assert.True(t, ok)

	// Check different routing ID
	ok = it.Verify("00000000", *pubKey)
	assert.False(t, ok)

	// Check different address in token
	it, err = NewInviteToken(addr, "12345678", time.Date(2010, 01, 04, 12, 34, 56, 0, time.UTC), *privKey)
	assert.NoError(t, err)
	ok = it.Verify("12345678", *pubKey)
	assert.True(t, ok)
	it.AddrHash = hash.New("doctor@evil!")
	ok = it.Verify("12345678", *pubKey)
	assert.False(t, ok)

	// Check different expiry in token
	it, err = NewInviteToken(addr, "12345678", time.Date(2010, 01, 04, 12, 34, 56, 0, time.UTC), *privKey)
	assert.NoError(t, err)
	ok = it.Verify("12345678", *pubKey)
	assert.True(t, ok)
	it.Expiry = time.Date(2099, 01, 04, 12, 34, 56, 0, time.UTC)
	ok = it.Verify("12345678", *pubKey)
	assert.False(t, ok)

	// Check token with different public key
	kt, err := bmcrypto.FindKeyType("rsa")
	assert.NoError(t, err)
	_, pubKey2, _ := bmcrypto.GenerateKeyPair(kt)
	it, err = NewInviteToken(addr, "12345678", time.Date(2010, 01, 04, 12, 34, 56, 0, time.UTC), *privKey)
	assert.NoError(t, err)
	ok = it.Verify("12345678", *pubKey2)
	assert.False(t, ok)
	ok = it.Verify("12345678", *pubKey)
	assert.True(t, ok)

	// Check if until time is checked
	timeNow = func() time.Time {
		return time.Date(2012, 12, 31, 12, 34, 56, 0, time.UTC)
	}
	it, err = NewInviteToken(addr, "12345678", time.Date(2010, 01, 04, 12, 34, 56, 0, time.UTC), *privKey)
	assert.NoError(t, err)
	ok = it.Verify("12345678", *pubKey)
	assert.False(t, ok)
}

func TestParseInviteToken(t *testing.T) {
	tok, err := ParseInviteToken("--4632632asdf325325252352")
	assert.Error(t, err)
	assert.Nil(t, tok)

	tok, err = ParseInviteToken("MTM0MmQxMGVkZGZiZGQ5OGM2YzRiOTNkZDY3NjA3M2M1MGM2YmVkOGY2MzgxOTA0NGVmOWM2YTZmZjM3MTk2NjoxMjM0NTY3ODoxMjYyNjA4NDk2Om4fFRenJUpabgJ+tfe/h++44PtmgpT8CMlhOr3F68JR0QnC6zMo+BkWXJtbPRMxFYYXQ4KQc1pT2S4PMKw5Cw4=")
	assert.NoError(t, err)
	assert.Equal(t, "12345678", tok.RoutingID)
	assert.Equal(t, int64(1262608496), tok.Expiry.Unix())
	assert.Equal(t, hash.Hash("1342d10eddfbdd98c6c4b93dd676073c50c6bed8f63819044ef9c6a6ff371966"), tok.AddrHash)
	assert.Equal(t, []byte{0x6e, 0x1f, 0x15, 0x17, 0xa7, 0x25, 0x4a, 0x5a, 0x6e, 0x2, 0x7e, 0xb5, 0xf7, 0xbf, 0x87, 0xef, 0xb8, 0xe0, 0xfb, 0x66, 0x82, 0x94, 0xfc, 0x8, 0xc9, 0x61, 0x3a, 0xbd, 0xc5, 0xeb, 0xc2, 0x51, 0xd1, 0x9, 0xc2, 0xeb, 0x33, 0x28, 0xf8, 0x19, 0x16, 0x5c, 0x9b, 0x5b, 0x3d, 0x13, 0x31, 0x15, 0x86, 0x17, 0x43, 0x82, 0x90, 0x73, 0x5a, 0x53, 0xd9, 0x2e, 0xf, 0x30, 0xac, 0x39, 0xb, 0xe}, tok.Signature)
}
