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

const expectedSignature = "gm8cy14MCib0V1vqi0YA02hyQbiEM+TYem/K1lJeu8N+A9WJIrjLtSsugpw9B5cpttc5wXCmp1Iau2JrA3SSkOo5ZKDieKEYpabGeMdnu4hfb2bvYs4T4Zm6m2/IaovWw/eF5zeeHF/0VUX5VQvQkkFUhfrse3VGdH0tRU9tpmwqfZL5MnKqMFp3a3e8ZVNLbUeB7tsopJfqiW8SdSptjgRmKUhvaya41Nn47jK6UmbVdWzTujqfsn0KxhQu6YV4qY5ItJf5WyIR+fbKqinUYCoDFinua+fL7j/nVeLcPEwLYoioXS+inTTTkrqpYLydjdS9x4lLr731suj9lAjNxQ=="

func Test_signHeader(t *testing.T) {
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

func Test_VerifyHeader(t *testing.T) {
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
	)

	ai = resolver.AddressInfo{
		Hash:        "000000000000000000000000000097026f0daeaec1aeb8351b096637679cf350",
		PublicKey:   *pubKey,
		RoutingID:   "87654321",
		Pow:         pow.String(),
		RoutingInfo: resolver.RoutingInfo{},
	}
	_ = repo.UploadAddress(&ai, *privKey, *pow)

	// Note: our sender uses key3
	privKey, pubKey, err = bmtest.ReadTestKey("../../testdata/key-3.json")
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

	return privKey
}
