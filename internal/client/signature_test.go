package client

import (
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"

	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	bmtest "github.com/bitmaelum/bitmaelum-suite/internal/testing"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/stretchr/testify/assert"
)

const expectedSignature = "HrjkBJQCzb5uTsToJ2FRd4bdJJgABwYP0GZl+JDhM0vIkOvqTvzTMex91wEjYbhhVcfHTGg+6TCY28Q5NhCNxGiCIo/5G/Nck3K9WLCeaN6eQWzyJgWmfQtp3J8rgDLFTsrZckV9LVRtHzV15rT1fP+2a4Y8oSMuWgPCiaOpckeB+/BYcdQcUmrVozt6Zk0cVIbVUMePVaaAXxgQWTVdeO1iw3zKWAaV0GQBJ6Y1OTLdbVgCz9f935koBxBkz+19pJXWzFWL32ECBTummAnJt7mar+0AkFLeLcvN+YiwNBJx/m8/5nptVLoSCY/12Cddinnt3TFT+XIHp1La5WrCfw=="

func TestSignHeader(t *testing.T) {
	privKey := setup()

	header := &message.Header{}
	_ = bmtest.ReadJSON("../../testdata/header-001.json", &header)
	assert.Empty(t, header.ClientSignature)
	err := SignHeader(header, *privKey)
	assert.NoError(t, err)
	assert.Equal(t, expectedSignature, header.ClientSignature)

	// Already present, don't overwrite
	_ = bmtest.ReadJSON("../../testdata/header-001.json", &header)
	assert.NotEmpty(t, header.ClientSignature)
	header.ClientSignature = "foobar"
	err = SignHeader(header, *privKey)
	assert.NoError(t, err)

	assert.Equal(t, "foobar", header.ClientSignature)
}

func TestVerifyHeader(t *testing.T) {
	privKey := setup()

	header := &message.Header{}
	_ = bmtest.ReadJSON("../../testdata/header-001.json", &header)
	assert.Empty(t, header.ClientSignature)
	err := SignHeader(header, *privKey)
	assert.NoError(t, err)
	assert.Equal(t, expectedSignature, header.ClientSignature)

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

func setup() *bmcrypto.PrivKey {
	// Setup container with mock repository for routing
	repo, _ := resolver.NewMockRepository()
	container.SetResolveService(resolver.KeyRetrievalService(repo))

	privKey, pubKey, err := bmtest.ReadTestKey("../../testdata/key-2.json")
	if err != nil {
		panic(err)
	}

	pow := proofofwork.NewWithoutProof(1, "foobar")
	var (
		ai resolver.AddressInfo
		ri resolver.RoutingInfo
	)

	privKey, pubKey, err = bmtest.ReadTestKey("../../testdata/key-1.json")
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
	privKey, pubKey, err = bmtest.ReadTestKey("../../testdata/key-3.json")
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
