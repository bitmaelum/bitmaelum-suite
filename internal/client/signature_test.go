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

package client

import (
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	testing2 "github.com/bitmaelum/bitmaelum-suite/internal/testing"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/stretchr/testify/assert"
)

const (
	expectedSignature  = "MsvxtxsxihfHRNDG20hmydSg7xjSIdaFmiz7iZVT+KLpvzP4isVimqzblN9bfxj+kLKuxKpXvSR6gS4Mi5yMwh1Xfsk0QJaVNJYNwB12i+kB5of2pm1vTipU4Y4I9E2Kj2TUfGuJZ51XeELPgr/d85ttrAk3DUKPXG5DcpEKqzJeocxryYq9968WB0Yal1SnxFc93HK1Lx5GWXGvXbVXnencCfZ+OYAVjuYOxtyAXQBEG2OjQUTfnmH+h8dyP0zIMl2xpku32uynp9Qy7gKyytACSFSSDTlHAAvZA8ARj7vrCGQWsUW0popMvx9N+z2KY1sL4q8/B/jg/uIpKa2AIQ=="
	expectedSignature2 = "Gy7Wz2Qj/kWBgO8r0ltYEbaS01EpNh0o2kyzmGsr7KPyH6X8Dl/uipYotEovbPEP5rhbip3Lwl6k1hO6FqqZE4v/2zyTdT+UN2o4Sm/BcsA00Cd7kttxyOdeqOpG+nJTLTi1nG/zbqSFli7qEAVtdUwMgs+w6/Vd0uv0SK4FnrcHqWIIRWBB0AbkMoVtR69J/zsPhvI6lA741cLM/P1K+CnrnK4OkoJ1nvXw9NX7AxBGh9yDyVMdAjLKjniLViXp+ZQwh9o/mIzus5mP8BS5q3tga9nhhg8k+IdXvs6kc0yURYhMESsbKroL+lLn1qvc0SdJkthjWhwa/e5MotXUCg=="
)

func TestSignHeader(t *testing.T) {
	privKey := setup()

	header := &message.Header{}
	_ = testing2.ReadJSON("../../testdata/header-001.json", &header)
	assert.Empty(t, header.ClientSignature)
	err := SignHeader(header, *privKey)
	assert.NoError(t, err)
	assert.Equal(t, expectedSignature2, header.ClientSignature)

	// Already present, don't overwrite
	_ = testing2.ReadJSON("../../testdata/header-001.json", &header)
	assert.NotEmpty(t, header.ClientSignature)
	header.ClientSignature = "foobar"
	err = SignHeader(header, *privKey)
	assert.NoError(t, err)

	assert.Equal(t, "foobar", header.ClientSignature)
}

func TestVerifyHeader(t *testing.T) {
	privKey := setup()

	header := &message.Header{}
	_ = testing2.ReadJSON("../../testdata/header-001.json", &header)
	assert.Empty(t, header.ClientSignature)
	err := SignHeader(header, *privKey)
	assert.NoError(t, err)
	assert.Equal(t, expectedSignature2, header.ClientSignature)

	// All is ok
	ok := VerifyHeader(*header)
	assert.True(t, ok)

	// Incorrect decoding
	header.ClientSignature = "A"
	ok = VerifyHeader(*header)
	assert.False(t, ok)

	// Empty sig is not ok
	header.ClientSignature = ""
	ok = VerifyHeader(*header)
	assert.False(t, ok)

	// incorrect key
	header.ClientSignature = "Zm9vYmFy"
	ok = VerifyHeader(*header)
	assert.False(t, ok)
}

func TestSignHeaderWithOnbehalfKey(t *testing.T) {
	_ = setup()
	// This is our onbehalf key
	privKey, _, _ := testing2.ReadTestKey("../../testdata/key-ed25519-2.json")

	header := &message.Header{}
	_ = testing2.ReadJSON("../../testdata/header-003.json", &header)
	assert.Empty(t, header.ClientSignature)
	err := SignHeader(header, *privKey)
	assert.NoError(t, err)
	assert.Equal(t, "ynMf3K4f0kZLIgcvw7I3cu9p9aBO8BlpAC8hswBQ3oUFtmpn+s6ND2iciDga7zEhP/hbWbGYFo8ITySpH3PIDA==", header.ClientSignature)

	ok := VerifyHeader(*header)
	assert.True(t, ok)
}

func setup() *bmcrypto.PrivKey {
	// Setup container with mock repository for routing
	repo, _ := resolver.NewMockRepository()
	container.SetResolveService(resolver.KeyRetrievalService(repo))

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
	_ = repo.UploadAddress(&ai, *privKey, *pow)

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
	_ = repo.UploadAddress(&ai, *privKey, *pow)

	return privKey
}
