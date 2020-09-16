package address

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ValidAddress(t *testing.T) {
	validAddresses := []string{
		"jay!",
		"jay@org!",
		"jay.@org!",
		"jay-@org!",
		"jay-@o-rg!",
		"jay-@o.rg!",
		"j1234!",
		"1ja!",
		"abc@de!",
		"yjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjay@yjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjay!",
	}

	for _, address := range validAddresses {
		assert.True(t, IsValidAddress(address), address)
	}
}

func Test_InvalidAddress(t *testing.T) {
	invalidAddresses := []string{
		"jay",
		"j!",
		"ja!",
		"1!",
		".@org!",
		"@@org!",
		"@org!",
		"ab@de!",
		"abc@d!",
		"jay",
		"yjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjay1@yjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjay1!",
		"yjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjay1@yjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjay!",
		"yjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjay@yjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjay1!",
		".jay!",
		"-jay!",
		"jay@-org!",
		"jay@.org!",
	}

	for _, address := range invalidAddresses {
		assert.False(t, IsValidAddress(address), address)
	}
}

func Test_Address(t *testing.T) {
	a, err := New("joshua@bitmaelum!")
	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.Equal(t, "joshua@bitmaelum!", a.String())

	a, err = New("joshua!")
	assert.NoError(t, err)
	assert.Equal(t, "joshua!", a.String())
	assert.Equal(t, HashAddress("66c94b6643ada5661b2d940eb87502b5af0f47f40fd45ce0fa125502dfa9c1ee"), a.Hash())

	a, err = New("j!")
	assert.Nil(t, a)
	assert.Error(t, err)

	a, err = New("JOHN@EXAMPLE!")
	assert.NotNil(t, a)
	assert.NoError(t, err)

	assert.Equal(t, "john", a.Local)
	assert.Equal(t, "example", a.Org)

	assert.Equal(t, "john@example!", a.String())
	assert.Equal(t, "16d0a463eb0be0514246e65b6b2d74c96d876bd1531f3bc095ac4b9f0b26d71c", a.Hash().String())

	a, err = New("JOHN!")
	assert.NoError(t, err)
	assert.Equal(t, "john!", a.String())

	a, err = New("JOHN@ex!")
	assert.NoError(t, err)
	a.Org = ""
	assert.Equal(t, "john!", a.String())
}

func Test_HashAddress(t *testing.T) {
	ha, err := NewHash("joshua@bitmaelum!")
	assert.NoError(t, err)
	assert.NotNil(t, ha)
	assert.Equal(t, "6b024a4e51c0c4a30c3750115c66be776253880bb4af0f313e3bf2236e808840", ha.String())

	ha, err = NewHash("incorrectaddress")
	assert.Error(t, err)
	assert.Nil(t, ha)

	ha, err = NewHashFromHash("6b024a4e51c0c4a30c3750115c66be776253880bb4af0f313e3bf2236e808840")
	assert.NoError(t, err)
	assert.NotNil(t, ha)
	assert.Equal(t, "6b024a4e51c0c4a30c3750115c66be776253880bb4af0f313e3bf2236e808840", ha.String())
}

func Test_Verify(t *testing.T) {
	assert.True(t, VerifyHash(
		"6b024a4e51c0c4a30c3750115c66be776253880bb4af0f313e3bf2236e808840",
		"fc52fabe94c0e037d2df4498e87481a6438960c9f73d517584a7a5c564535ac4",
		"49aa67181f4a3176f9b65605390bb81126e8ff1f6d03b1bd220c53e7a6b36d3e",
	))

	assert.False(t, VerifyHash(
		"6b024a4e51c0c4a30c3750115c66be776253880bb4af0f313e3bf2236e808840",
		"fc52fabe94c0e037d2df4498e87481a6438960c9f73d517584a7a5c564535ac4",
		"0000000000000006f9b65605390bb81126e8ff1f6d03b1bd220c53e7a6b36d3e",
	))

}
