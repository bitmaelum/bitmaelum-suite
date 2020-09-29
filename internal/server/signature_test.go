package server

import (
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_signHeader(t *testing.T) {
	setup()

	header := &message.Header{}
	_ = internal.ReadHeader("../../testdata/header-001.json", &header)
	assert.Empty(t, header.ServerSignature)
	err := SignHeader(header)
	assert.NoError(t, err)

	assert.Equal(t, "lNSBxGYk7Gn9laP+jA1GSDd+KJxz1hx26MCJSPBkV5tMZ3QGxIKSkpAMc2fmZfnCm4dZrGOiorNuTdJAgT87FodVH0Q7OiU7lMJq8FpxgbMrQq2FIp0PKoQ3XfVvrj3S6pY23W5aMMRa8Qn2y67PSu4KZsSYJzjMU1PHEQMUCoVTsM/a/5FmCFnTy5e7FpnU5ssEn9joYQWFGvu1kYUhJwuqfmr467acJBvyN3OME0t3M5csRZucOeSZBU9TvMrt9dv85ATvUthyM/aTEnxEqemTJDfag/+5gyUmSUtDnGwWQkGBESkeUC7YaXIGpcMutzJmvSo5RhQ8KuAG8jjo9g==", header.ServerSignature)

	// Already present, don't overwrite
	_ = internal.ReadHeader("../../testdata/header-001.json", &header)
	assert.NotEmpty(t, header.ServerSignature)
	header.ServerSignature = "foobar"
	err = SignHeader(header)
	assert.NoError(t, err)

	assert.Equal(t, "foobar", header.ServerSignature)
}

func Test_VerifyHeader(t *testing.T) {
	setup()

	header := &message.Header{}
	_ = internal.ReadHeader("../../testdata/header-001.json", &header)
	assert.Empty(t, header.ServerSignature)
	err := SignHeader(header)
	assert.NoError(t, err)
	assert.Equal(t, "lNSBxGYk7Gn9laP+jA1GSDd+KJxz1hx26MCJSPBkV5tMZ3QGxIKSkpAMc2fmZfnCm4dZrGOiorNuTdJAgT87FodVH0Q7OiU7lMJq8FpxgbMrQq2FIp0PKoQ3XfVvrj3S6pY23W5aMMRa8Qn2y67PSu4KZsSYJzjMU1PHEQMUCoVTsM/a/5FmCFnTy5e7FpnU5ssEn9joYQWFGvu1kYUhJwuqfmr467acJBvyN3OME0t3M5csRZucOeSZBU9TvMrt9dv85ATvUthyM/aTEnxEqemTJDfag/+5gyUmSUtDnGwWQkGBESkeUC7YaXIGpcMutzJmvSo5RhQ8KuAG8jjo9g==", header.ServerSignature)

	// All is ok
	ok := VerifyHeader(*header)
	assert.True(t, ok)

	// Incorect decoding
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
	privKey, pubKey, err := internal.ReadTestKey("../../testdata/key-1.json")
	if err != nil {
		panic(err)
	}
	config.Server.Routing = &config.Routing{
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

	privKey, pubKey, err = internal.ReadTestKey("../../testdata/key-2.json")
	if err != nil {
		panic(err)
	}
	ai = resolver.AddressInfo{
		Hash:        "000000000000000000000000000097026f0daeaec1aeb8351b096637679cf350",
		PublicKey:   *pubKey,
		RoutingID:   "87654321",
		Pow:         pow.String(),
		RoutingInfo: resolver.RoutingInfo{},
	}
	_ = repo.UploadAddress(&ai, *privKey, *pow)

	privKey, pubKey, err = internal.ReadTestKey("../../testdata/key-3.json")
	if err != nil {
		panic(err)
	}
	ai = resolver.AddressInfo{
		Hash:        "000000000000000000018f66a0f3591a883f2b9cc3e95a497e7cf9da1071b4cc",
		PublicKey:   *pubKey,
		RoutingID:   "12345678",
		Pow:         pow.String(),
		RoutingInfo: resolver.RoutingInfo{},
	}
	_ = repo.UploadAddress(&ai, *privKey, *pow)

	// Note: our mail server uses key1
	privKey, pubKey, err = internal.ReadTestKey("../../testdata/key-1.json")
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
