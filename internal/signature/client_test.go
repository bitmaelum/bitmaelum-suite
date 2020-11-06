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

package signature

import (
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	testing2 "github.com/bitmaelum/bitmaelum-suite/internal/testing"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/stretchr/testify/assert"
)

const (
	expectedClientSignature = "lIqI1QYBRHl7yRW367Lx2n/PFadrYDZ2a2NGSaL40EKum0ncOIXs8CIqKZ+LCUgmK2a9iH2d3mbXVPwZ3PBGsVgReaomyG6NrDbZ0PCbgnjmrmkVAFV0bDHlOxUl/BzyV+seIL7FL0lu+cODaHkmzH16FsZ5Vqcf1/Qe2GR/0Ka6xbWcIcajGsKtTx+WtGeZGZ5oLbAFatEjiv5gMAn2umKpP+w7uKhPa6CsYkv2YMVw+z/1NU2CO0jE6/2muihF9x4nPw6yiy+sXP86B26FQXLBcMgTZ4TAtzr/b2KvcEDj8y8HISs/YHJvTdqAXzYTPnha37ZIIZ7ce27Z41GAUQ=="
)

func TestSignClientHeader(t *testing.T) {
	privKey := setupClient()

	header := &message.Header{}
	_ = testing2.ReadJSON("../../testdata/header-001.json", &header)
	assert.Empty(t, header.Signatures.Client)
	err := SignClientHeader(header, *privKey)
	assert.NoError(t, err)
	assert.Equal(t, expectedSignature, header.Signatures.Client)

	// Already present, don't overwrite
	_ = testing2.ReadJSON("../../testdata/header-001.json", &header)
	assert.NotEmpty(t, header.Signatures.Client)
	header.Signatures.Client = "foobar"
	err = SignClientHeader(header, *privKey)
	assert.NoError(t, err)

	assert.Equal(t, "foobar", header.Signatures.Client)
}

func TestVerifyClientHeader(t *testing.T) {
	privKey := setupClient()

	header := &message.Header{}
	_ = testing2.ReadJSON("../../testdata/header-001.json", &header)
	assert.Empty(t, header.Signatures.Client)
	err := SignClientHeader(header, *privKey)
	assert.NoError(t, err)
	assert.Equal(t, expectedSignature, header.Signatures.Client)

	// All is ok
	ok := VerifyClientHeader(*header)
	assert.True(t, ok)

	// Incorrect decoding
	header.Signatures.Client = "A"
	ok = VerifyClientHeader(*header)
	assert.False(t, ok)

	// Empty sig is not ok
	header.Signatures.Client = ""
	ok = VerifyClientHeader(*header)
	assert.False(t, ok)

	// incorrect key
	header.Signatures.Client = "Zm9vYmFy"
	ok = VerifyClientHeader(*header)
	assert.False(t, ok)
}

func setupClient() *bmcrypto.PrivKey {
	addr, _ := address.NewAddress("foobar!")

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

	privKey, pubKey, err := testing2.ReadTestKey("../../testdata/key-1.json")
	if err != nil {
		panic(err)
	}
	ai = resolver.AddressInfo{
		RoutingID:   "87654321",
		PublicKey:   *pubKey,
		RoutingInfo: resolver.RoutingInfo{},
		Pow:         pow.String(),
		Hash:        "000000000000000000000000000097026f0daeaec1aeb8351b096637679cf350",
	}
	_ = repo.UploadAddress(*addr, &ai, *privKey, *pow, "")

	ri = resolver.RoutingInfo{
		PublicKey: *pubKey,
		Routing:   "127.0.0.1",
		Hash:      "12345678",
	}

	_ = repo.UploadRouting(&ri, *privKey)

	// Note: our sender uses key3
	privKey, pubKey, err = testing2.ReadTestKey("../../testdata/key-3.json")
	if err != nil {
		panic(err)
	}

	ai = resolver.AddressInfo{
		RoutingID:   "12345678",
		RoutingInfo: resolver.RoutingInfo{},
		PublicKey:   *pubKey,
		Hash:        "000000000000000000018f66a0f3591a883f2b9cc3e95a497e7cf9da1071b4cc",
		Pow:         pow.String(),
	}
	_ = repo.UploadAddress(*addr, &ai, *privKey, *pow, "")

	return privKey
}
