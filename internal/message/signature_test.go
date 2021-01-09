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

package message

import (
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	testing2 "github.com/bitmaelum/bitmaelum-suite/internal/testing"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/stretchr/testify/assert"
)

const (
	expectedServerSignature = "h3+FwOJV+FuO2fq42X4XiRkn8Tri2n2orN2ELRKc0E7484HzmpIhIDptGd2m84zsEyKWfi5yoL5k5Ld7p44Odvuh2qtE9S089HONHtH5rcy5kkPoUO0zCCpSrrNPBxuIzfpge64aCihSMHuxJkUoF40iHX0ZG7VNqm1K1RRsgd9e5v1chsr7MepMLCw7Xh4Vc3WTdIV98dBpfUh3XsGNxnIB6xXMPh/mexi35v8MUjgSICEn1Oa9h15WmLdbMictEkWa+yjySM9Udv5rSLSU3lRUzmj6RIqlsKYRegkOAQsXep2d58KFgD86at0ylBBdyeghnQO5aaRuoHpHWbEfjQ=="
	expectedClientSignature = "qA7PuGm+wTv4bsdWdjaadLeE4urqgV09c1vGnS1ashP54pCSpT/melrGkf7hP36CupnVXMlGvHnC4jEQZM5BJvuThimB0MRYWQta2/2s8xlXfe/crYnyoIEQ58vV9iD/LXQhFxMEQz9K8KO1HuLie3IJYs48yhwrlPn0uKih2BD5mOz6a+bZSex77Karx5N7lcC5Uz4Xr2CV6pTaO2achN2NnTbCuO1HYzx5n5DNe6HfFRKlTIpIoHbqnJOi8f6/wXlCGkpT1+GYvlED+6fowJT7fmXfxePYxB36vMJjhS1P5d1c1bcdYR0ve2kiSNehpF9ScY+A49Ovn3fqRuXL1g=="
)

func TestSignServerHeader(t *testing.T) {
	setupServer()

	header := &Header{}
	_ = testing2.ReadJSON("../../testdata/header-002.json", &header)
	assert.Empty(t, header.Signatures.Server)
	err := SignServerHeader(header)
	assert.NoError(t, err)

	assert.Equal(t, expectedServerSignature, header.Signatures.Server)

	// // Already present, don't overwrite
	// _ = testing2.ReadJSON("../../testdata/header-002.json", &header)
	// assert.NotEmpty(t, header.Signatures.Server)
	// header.Signatures.Server = "foobar"
	// err = SignServerHeader(header)
	// assert.NoError(t, err)
	//
	// assert.Equal(t, "foobar", header.Signatures.Server)
}

func TestVerifyServerHeader(t *testing.T) {
	setupServer()

	header := &Header{}
	_ = testing2.ReadJSON("../../testdata/header-002.json", &header)
	assert.Empty(t, header.Signatures.Server)
	err := SignServerHeader(header)
	assert.NoError(t, err)
	assert.Equal(t, expectedServerSignature, header.Signatures.Server)

	// All is ok
	ok := VerifyServerHeader(*header)
	assert.True(t, ok)

	// Unknown addr
	addr := header.From.Addr
	header.From.Addr = "fooobar"
	ok = VerifyServerHeader(*header)
	assert.False(t, ok)
	header.From.Addr = addr

	// Incorrect decoding
	header.Signatures.Server = "A"
	ok = VerifyServerHeader(*header)
	assert.False(t, ok)

	// Empty sig is not ok
	header.Signatures.Server = ""
	ok = VerifyServerHeader(*header)
	assert.False(t, ok)

	// incorrect key
	header.Signatures.Server = "Zm9vYmFy"
	ok = VerifyServerHeader(*header)
	assert.False(t, ok)

	// Test for true if the message is server-side
	header.From.SignedBy = SignedByTypeServer
	header.Signatures.Server = "Zm9vYmFy"
	ok = VerifyServerHeader(*header)
	assert.True(t, ok)

}

func TestSignClientHeader(t *testing.T) {
	privKey := setupClient()

	header := &Header{}
	_ = testing2.ReadJSON("../../testdata/header-001.json", &header)
	assert.Empty(t, header.Signatures.Client)
	err := SignClientHeader(header, *privKey)
	assert.NoError(t, err)
	assert.Equal(t, expectedClientSignature, header.Signatures.Client)

	// Already present, overwrite with the correct
	_ = testing2.ReadJSON("../../testdata/header-001.json", &header)
	assert.NotEmpty(t, header.Signatures.Client)
	header.Signatures.Client = "foobar"
	err = SignClientHeader(header, *privKey)
	assert.NoError(t, err)
	assert.Equal(t, expectedClientSignature, header.Signatures.Client)
}

func TestVerifyClientHeader(t *testing.T) {
	privKey := setupClient()

	header := &Header{}
	_ = testing2.ReadJSON("../../testdata/header-001.json", &header)
	assert.Empty(t, header.Signatures.Client)
	err := SignClientHeader(header, *privKey)
	assert.NoError(t, err)
	assert.Equal(t, expectedClientSignature, header.Signatures.Client)

	// All is ok
	ok := VerifyClientHeader(*header)
	assert.True(t, ok)

	// Unknown addr
	addr := header.From.Addr
	header.From.Addr = "5723579275927597239572935729"
	ok = VerifyClientHeader(*header)
	assert.False(t, ok)
	header.From.Addr = addr

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

func TestVerifyClientHeaderWithServerSignature(t *testing.T) {
	_ = setupClient()

	header := &Header{}
	_ = testing2.ReadJSON("../../testdata/header-001.json", &header)

	header.From.SignedBy = SignedByTypeServer

	// Test with correct routing ID
	header.From.Addr = "12345678"
	header.Signatures.Client = "6i4S/htfh8Ye2e3higYHkQZ5KpIejj6U7sXPLwU2CwAEDYNVZ2OvrFjbRjq49RmySYn8sDvnkhpLr1WRr78+Cg=="
	ok := VerifyClientHeader(*header)
	assert.True(t, ok)

	// Correct sig, wrong routing ID
	header.From.Addr = "44444444"
	header.Signatures.Client = "6i4S/htfh8Ye2e3higYHkQZ5KpIejj6U7sXPLwU2CwAEDYNVZ2OvrFjbRjq49RmySYn8sDvnkhpLr1WRr78+Cg=="
	ok = VerifyClientHeader(*header)
	assert.False(t, ok)

	// Correct routing ID, no sig
	header.From.Addr = "12345678"
	header.Signatures.Client = ""
	ok = VerifyClientHeader(*header)
	assert.False(t, ok)

	// empty from addr
	header.From.Addr = ""
	header.Signatures.Client = "6i4S/htfh8Ye2e3higYHkQZ5KpIejj6U7sXPLwU2CwAEDYNVZ2OvrFjbRjq49RmySYn8sDvnkhpLr1WRr78+Cg=="
	ok = VerifyClientHeader(*header)
	assert.False(t, ok)

	// Wrong from addr (not a routing)
	header.From.Addr = "000000000000000000018f66a0f3591a883f2b9cc3e95a497e7cf9da1071b4cc"
	header.Signatures.Client = "6i4S/htfh8Ye2e3higYHkQZ5KpIejj6U7sXPLwU2CwAEDYNVZ2OvrFjbRjq49RmySYn8sDvnkhpLr1WRr78+Cg=="
	ok = VerifyClientHeader(*header)
	assert.False(t, ok)
}

func TestVerifyClientHeaderWithOnbehalfSignature(t *testing.T) {
	_ = setupClient()

	header := &Header{}
	_ = testing2.ReadJSON("../../testdata/header-003.json", &header)

	ok := VerifyClientHeader(*header)
	assert.True(t, ok)

	// Incorrect authorizer signature
	header.AuthorizedBy.Signature = "Zm9vYmFy"
	ok = VerifyClientHeader(*header)
	assert.False(t, ok)

	// Incorrect encoded signature
	header.AuthorizedBy.Signature = "foobar"
	ok = VerifyClientHeader(*header)
	assert.False(t, ok)
}

func setupClient() *bmcrypto.PrivKey {
	addr, _ := address.NewAddress("foobar!")

	// Setup container with mock repository for routing
	repo, _ := resolver.NewMockRepository()
	container.Instance.SetShared("resolver", func() (interface{}, error) {
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

	// Note: our sender uses key3
	privKey, pubKey, err = testing2.ReadTestKey("../../testdata/key-ed25519-3.json")
	if err != nil {
		panic(err)
	}
	ri = resolver.RoutingInfo{
		PublicKey: *pubKey,
		Routing:   "127.1.2.3",
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

func setupServer() {
	addr, _ := address.NewAddress("foobar!")

	// Note: our mail server uses key1
	kp, err := testing2.ReadKeyPair("../../testdata/key-1.json")
	if err != nil {
		panic(err)
	}
	config.Routing = config.RoutingConfig{
		RoutingID: "12345678",
		KeyPair:   kp,
	}

	// Setup container with mock repository for routing
	repo, _ := resolver.NewMockRepository()
	container.Instance.SetShared("resolver", func() (interface{}, error) {
		return resolver.KeyRetrievalService(repo), nil
	})

	uploadAddress(repo, *addr, "111100000000000000000000000097026f0daeaec1aeb8351b096637679cf350", "87654321", "../../testdata/key-2.json")
	uploadAddress(repo, *addr, "111100000000000000018f66a0f3591a883f2b9cc3e95a497e7cf9da1071b4cc", "12345678", "../../testdata/key-3.json")

	// Note: our mail server uses key1
	privKey, pubKey, err := testing2.ReadTestKey("../../testdata/key-1.json")
	if err != nil {
		panic(err)
	}
	ri := resolver.RoutingInfo{
		Hash:      "12345678",
		PublicKey: *pubKey,
		Routing:   "127.0.0.1",
	}
	_ = repo.UploadRouting(&ri, *privKey)
}

func uploadAddress(repo resolver.AddressRepository, addr address.Address, addrHash string, routingID string, keyPath string) {
	pow := proofofwork.NewWithoutProof(1, "foobar")

	privKey, pubKey, err := testing2.ReadTestKey(keyPath)
	if err != nil {
		panic(err)
	}

	ai := resolver.AddressInfo{
		Hash:        addrHash,
		PublicKey:   *pubKey,
		RoutingID:   routingID,
		Pow:         pow.String(),
		RoutingInfo: resolver.RoutingInfo{},
	}
	_ = repo.UploadAddress(addr, &ai, *privKey, *pow, "")

}
