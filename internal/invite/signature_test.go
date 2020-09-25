package invite

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSignature(t *testing.T) {
	addr, _ := address.NewHash("john@acme!")

	privKey, err := bmcrypto.NewPrivKey("ed25519 MC4CAQAwBQYDK2VwBCIEILq+V/CUlMdbmoQC1odEgOEmtMBQu0UpIICxJbQM1vhd")
	assert.NoError(t, err)
	pubKey, err := bmcrypto.NewPubKey("ed25519 MCowBQYDK2VwAyEARdZSwluYtMWTGI6Rvl0Bhu40RBDn6D88wyzFL1IR3DU=")
	assert.NoError(t, err)

	// Assume this is the current time during tests
	timeNow = func() time.Time {
		return time.Date(2010, 01, 01, 12, 34, 56, 0, time.UTC)
	}

	it, err := NewInviteToken(*addr, "12345678", time.Date(2010, 01, 04, 12, 34, 56, 0, time.UTC), *privKey)
	assert.NoError(t, err)
	assert.Equal(t, "MTM0MmQxMGVkZGZiZGQ5OGM2YzRiOTNkZDY3NjA3M2M1MGM2YmVkOGY2MzgxOTA0NGVmOWM2YTZmZjM3MTk2NjoxMjM0NTY3ODoxMjYyNjA4NDk2Om4fFRenJUpabgJ+tfe/h++44PtmgpT8CMlhOr3F68JR0QnC6zMo+BkWXJtbPRMxFYYXQ4KQc1pT2S4PMKw5Cw4=", it.String())
	ok := it.Verify("12345678", *pubKey)
	assert.True(t, ok)

	// Check different routing ID
	ok = it.Verify("00000000", *pubKey)
	assert.False(t, ok)

	// Check different address in token
	it, err = NewInviteToken(*addr, "12345678", time.Date(2010, 01, 04, 12, 34, 56, 0, time.UTC), *privKey)
	assert.NoError(t, err)
	ok = it.Verify("12345678", *pubKey)
	assert.True(t, ok)
	addr2, _ := address.NewHash("doctor@evil!")
	it.Address = *addr2
	ok = it.Verify("12345678", *pubKey)
	assert.False(t, ok)

	// Check different expiry in token
	it, err = NewInviteToken(*addr, "12345678", time.Date(2010, 01, 04, 12, 34, 56, 0, time.UTC), *privKey)
	assert.NoError(t, err)
	ok = it.Verify("12345678", *pubKey)
	assert.True(t, ok)
	it.Expiry = time.Date(2099, 01, 04, 12, 34, 56, 0, time.UTC)
	ok = it.Verify("12345678", *pubKey)
	assert.False(t, ok)

	// Check token with differnet public key
	_, pubKey2, _ := encrypt.GenerateKeyPair(bmcrypto.KeyTypeRSA)
	it, err = NewInviteToken(*addr, "12345678", time.Date(2010, 01, 04, 12, 34, 56, 0, time.UTC), *privKey)
	assert.NoError(t, err)
	ok = it.Verify("12345678", *pubKey2)
	assert.False(t, ok)
	ok = it.Verify("12345678", *pubKey)
	assert.True(t, ok)

	// Check if until time is checked
	timeNow = func() time.Time {
		return time.Date(2012, 12, 31, 12, 34, 56, 0, time.UTC)
	}
	it, err = NewInviteToken(*addr, "12345678", time.Date(2010, 01, 04, 12, 34, 56, 0, time.UTC), *privKey)
	assert.NoError(t, err)
	ok = it.Verify("12345678", *pubKey)
	assert.False(t, ok)
}

/*

func TestSignature(t *testing.T) {
	addr, _ := address.NewHash("john@acme!")

	privKey, err := bmcrypto.NewPrivKey("ed25519 MC4CAQAwBQYDK2VwBCIEILq+V/CUlMdbmoQC1odEgOEmtMBQu0UpIICxJbQM1vhd")
	assert.NoError(t, err)
	pubKey, err := bmcrypto.NewPubKey("ed25519 MCowBQYDK2VwAyEARdZSwluYtMWTGI6Rvl0Bhu40RBDn6D88wyzFL1IR3DU=")
	assert.NoError(t, err)

	timeNow = func() time.Time {
		return time.Date(2010, 01, 01, 12, 34, 56, 0, time.UTC)
	}
	// Date must be in the future of our timenow func
	until := time.Date(2010, 01, 04, 12, 34, 56, 0, time.UTC)
	token, err := GenerateToken(*addr, "12345678", until, *privKey)
	assert.NoError(t, err)
	assert.Equal(t, "MTM0MmQxMGVkZGZiZGQ5OGM2YzRiOTNkZDY3NjA3M2M1MGM2YmVkOGY2MzgxOTA0NGVmOWM2YTZmZjM3MTk2NjoxMjM0NTY3ODoxMjYyNjA4NDk2Om4fFRenJUpabgJ+tfe/h++44PtmgpT8CMlhOr3F68JR0QnC6zMo+BkWXJtbPRMxFYYXQ4KQc1pT2S4PMKw5Cw4=", token)
	ok := VerifyToken(token, *pubKey)
	assert.True(t, ok)

	// Check different routing ID
	token2 := "MTM0MmQxMGVkZGZiZGQ5OGM2YzRiOTNkZDY3NjA3M2M1MGM2YmVkOGY2MzgxOTA0NGVmOWM2YTZmZjM3MTk2NjowMDAwMDAwMDoxMjYyNjA4NDk2Om4fFRenJUpabgJ+tfe/h++44PtmgpT8CMlhOr3F68JR0QnC6zMo+BkWXJtbPRMxFYYXQ4KQc1pT2S4PMKw5Cw4="
	ok = VerifyToken(token2, *pubKey)
	assert.False(t, ok)

	// Check different address in token
	token3 := "MDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDA6MTIzNDU2Nzg6MTI2MjYwODQ5NjpuHxUXpyVKWm4CfrX3v4fvuOD7ZoKU/AjJYTq9xevCUdEJwuszKPgZFlybWz0TMRWGF0OCkHNaU9kuDzCsOQsO"
	ok = VerifyToken(token3, *pubKey)
	assert.False(t, ok)

	// Check different expiry in token
	token4 := "MTM0MmQxMGVkZGZiZGQ5OGM2YzRiOTNkZDY3NjA3M2M1MGM2YmVkOGY2MzgxOTA0NGVmOWM2YTZmZjM3MTk2NjoxMjM0NTY3ODoxMDAxMTExMTExOm4fFRenJUpabgJ+tfe/h++44PtmgpT8CMlhOr3F68JR0QnC6zMo+BkWXJtbPRMxFYYXQ4KQc1pT2S4PMKw5Cw4="
	ok = VerifyToken(token4, *pubKey)
	assert.False(t, ok)


	// Check token with differnet public key
	_, pubKey2, err := encrypt.GenerateKeyPair(bmcrypto.KeyTypeRSA)
	assert.NoError(t, err)
	ok = VerifyToken(token, *pubKey2)
	assert.False(t, ok)


	// Check if until time is checked
	timeNow = func() time.Time {
		return time.Date(2012, 12, 31, 12, 34, 56, 0, time.UTC)
	}
	ok = VerifyToken(token, *pubKey)
	assert.False(t, ok)
}
*/
