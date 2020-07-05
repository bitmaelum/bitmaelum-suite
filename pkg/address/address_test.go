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
	assert.Equal(t, HashAddress("2f92571b5567b4557b94ac5701fc48e552ba9970d6dac89f7c2ebce92f1cd836"), a.Hash())

	a, err = New("j!")
	assert.Nil(t, a)
	assert.Error(t, err)

	a, err = New("JOHN@EXAMPLE!")
	assert.NotNil(t, a)
	assert.NoError(t, err)

	assert.Equal(t, "john", a.Local)
	assert.Equal(t, "example", a.Org)

	assert.Equal(t, "john@example!", a.String())
	assert.Equal(t, "f454fe8d4b5017369f9e64861f0d471efe3cdcbdf45732f26b7a377c3e93d278", a.Hash().String())

	a, err = New("JOHN!")
	assert.NoError(t, err)
	assert.Equal(t, "john!", a.String())

	a, err = New("JOHN@ex!")
	assert.NoError(t, err)
	a.Org = ""
	assert.Equal(t, "john!", a.String())

}
