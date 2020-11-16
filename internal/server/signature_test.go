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
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/stretchr/testify/assert"
)

const (
	expectedSignature = "smq7qWrkwTL1vquhEAn/WoRDZBT1BjUQaUSsmTSSPePLRpM1sjX10mJwxXlivYOIlgTtJ+0SIeBM9rDPBSTw5JJhOS9ZFmpAYEPG9tkU+9EjxfBNnEBPrsYaqE81tO7OtY4xFrlYecdhepGUQSxbZQU4+Ih5jE9jLb+SuUxR6lGw2u2P+Ngy75dD33zlMTTgnmaVTxBueRlfArDExW5QE+pFv9/uFi8xM5a7eGnQHSjufQ9gM6WOWzhyWAnKI+6XMwx3MoQ53H3OU2vn4tPSUQQyxB+L4WTH9JtC0nLC0ggzvo5LOdCw4rCljsciYiEZ2WssGD9kXLGFIU/ixEX2Kw=="
)

func TestSignHeader(t *testing.T) {
	setup()

	header := &message.Header{}
	_ = bmtest.ReadJSON("../../testdata/header-002.json", &header)
	assert.Empty(t, header.Signatures.Server)
	err := SignHeader(header)
	assert.NoError(t, err)

	assert.Equal(t, expectedSignature, header.Signatures.Server)

	// Already present, don't overwrite
	_ = bmtest.ReadJSON("../../testdata/header-002.json", &header)
	assert.NotEmpty(t, header.Signatures.Server)
	header.Signatures.Server = "foobar"
	err = SignHeader(header)
	assert.NoError(t, err)

	assert.Equal(t, "foobar", header.Signatures.Server)
}

func TestVerifyHeader(t *testing.T) {
	setup()

	header := &message.Header{}
	_ = bmtest.ReadJSON("../../testdata/header-002.json", &header)
	assert.Empty(t, header.Signatures.Server)
	err := SignHeader(header)
	assert.NoError(t, err)
	assert.Equal(t, expectedSignature, header.Signatures.Server)

	// All is ok
	ok := VerifyHeader(*header)
	assert.True(t, ok)

	// Incorrect decoding
	header.Signatures.Server = "A"
	ok = VerifyHeader(*header)
	assert.False(t, ok)

	// Empty sig is not ok
	header.Signatures.Server = ""
	ok = VerifyHeader(*header)
	assert.False(t, ok)

	// incorrect key
	header.Signatures.Server = "Zm9vYmFy"
	ok = VerifyHeader(*header)
	assert.False(t, ok)
}

func setup() {
	addr, _ := address.NewAddress("foobar!")

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
	container.Set("resolver", func() (interface{}, error) {
		return resolver.KeyRetrievalService(repo), nil
	})

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
	_ = repo.UploadAddress(*addr, &ai, *privKey, *pow, "")

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
	_ = repo.UploadAddress(*addr, &ai, *privKey, *pow, "")

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
