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

package server

import (
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	bmtest "github.com/bitmaelum/bitmaelum-suite/internal/testing"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/stretchr/testify/assert"
)

func TestSignHeader(t *testing.T) {
	setup()

	header := &message.Header{}
	_ = bmtest.ReadJSON("../../testdata/header-002.json", &header)
	assert.Empty(t, header.ServerSignature)
	err := SignHeader(header)
	assert.NoError(t, err)

	assert.Equal(t, "r4u8tkGcb9Wq2k3BSZQ5P2cipKBKH4VB524EEG/3Ijul6I1EuyR5RkgJoJKPjkSlN4XhLecnG/M0DBTcTA3H5JAitprZUTA0Va1rYktFBP/hzMKR/WehkaAshGm4+L5bf7JIpo8EzGU7CbWkyhwSn5GhlGN6deWRvpdy0fWXPRZqbdzAoEtVqdJKOAF1JkMMsXQfWeFzY1pcIYDqVSPOM0JFYRmcAyLaismgy42UV9QEy7Teyvy38aBraMylwFYJwwMKwoZjc9fDwcmMptWWKeLorksn3Jj86WXi6WCsOq5/BK3comjInfO8BKJ3pDDy8YUwOaUZuSTrlIfxXJvxxA==", header.ServerSignature)

	// Already present, don't overwrite
	_ = bmtest.ReadJSON("../../testdata/header-002.json", &header)
	assert.NotEmpty(t, header.ServerSignature)
	header.ServerSignature = "foobar"
	err = SignHeader(header)
	assert.NoError(t, err)

	assert.Equal(t, "foobar", header.ServerSignature)
}

func TestVerifyHeader(t *testing.T) {
	setup()

	header := &message.Header{}
	_ = bmtest.ReadJSON("../../testdata/header-002.json", &header)
	assert.Empty(t, header.ServerSignature)
	err := SignHeader(header)
	assert.NoError(t, err)
	assert.Equal(t, "r4u8tkGcb9Wq2k3BSZQ5P2cipKBKH4VB524EEG/3Ijul6I1EuyR5RkgJoJKPjkSlN4XhLecnG/M0DBTcTA3H5JAitprZUTA0Va1rYktFBP/hzMKR/WehkaAshGm4+L5bf7JIpo8EzGU7CbWkyhwSn5GhlGN6deWRvpdy0fWXPRZqbdzAoEtVqdJKOAF1JkMMsXQfWeFzY1pcIYDqVSPOM0JFYRmcAyLaismgy42UV9QEy7Teyvy38aBraMylwFYJwwMKwoZjc9fDwcmMptWWKeLorksn3Jj86WXi6WCsOq5/BK3comjInfO8BKJ3pDDy8YUwOaUZuSTrlIfxXJvxxA==", header.ServerSignature)

	// All is ok
	ok := VerifyHeader(*header)
	assert.True(t, ok)

	// Incorrect decoding
	header.ServerSignature = "A"
	ok = VerifyHeader(*header)
	assert.False(t, ok)

	// Empty sig is not ok
	header.ServerSignature = ""
	ok = VerifyHeader(*header)
	assert.False(t, ok)

	// incorrect key
	header.ServerSignature = "Zm9vYmFy"
	ok = VerifyHeader(*header)
	assert.False(t, ok)
}

func setup() {
	// Note: our mail server uses key1
	privKey, pubKey, err := bmtest.ReadTestKey("../../testdata/key-1.json")
	if err != nil {
		panic(err)
	}
	config.Routing = config.RoutingConfig{
		RoutingID:  "12345678",
		PrivateKey: *privKey,
		PublicKey:  *pubKey,
	}

	// Setup container with mock repository for routing
	repo, _ := resolver.NewMockRepository()
	container.SetResolveService(resolver.KeyRetrievalService(repo))

	pow := proofofwork.NewWithoutProof(1, "foobar")
	var (
		ai resolver.AddressInfo
		ri resolver.RoutingInfo
	)

	privKey, pubKey, err = bmtest.ReadTestKey("../../testdata/key-2.json")
	if err != nil {
		panic(err)
	}
	ai = resolver.AddressInfo{
		Hash:        "111100000000000000000000000097026f0daeaec1aeb8351b096637679cf350",
		PublicKey:   *pubKey,
		RoutingID:   "87654321",
		Pow:         pow.String(),
		RoutingInfo: resolver.RoutingInfo{},
	}
	_ = repo.UploadAddress(&ai, *privKey, *pow)

	privKey, pubKey, err = bmtest.ReadTestKey("../../testdata/key-3.json")
	if err != nil {
		panic(err)
	}
	ai = resolver.AddressInfo{
		Hash:        "111100000000000000018f66a0f3591a883f2b9cc3e95a497e7cf9da1071b4cc",
		PublicKey:   *pubKey,
		RoutingID:   "12345678",
		Pow:         pow.String(),
		RoutingInfo: resolver.RoutingInfo{},
	}
	_ = repo.UploadAddress(&ai, *privKey, *pow)

	// Note: our mail server uses key1
	privKey, pubKey, err = bmtest.ReadTestKey("../../testdata/key-1.json")
	if err != nil {
		panic(err)
	}
	ri = resolver.RoutingInfo{
		Hash:      "12345678",
		PublicKey: *pubKey,
		Routing:   "127.0.0.1",
	}
	_ = repo.UploadRouting(&ri, *privKey)
}
