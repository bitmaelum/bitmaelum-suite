package address

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidAddress(t *testing.T) {
	validAddresses := []string{
		"jay!",
		"jay@org!",
		"ja.y@org!",
		"j.a.y@org!",
		"j.a.y@o..rg!",
		"j.a.y@o..r---g!",
		"j1234!",
		"1ja!",
		"abc@de!",
		"yjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjay!",
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
		"jay.@org!",
		"jay-@org!",
		"jay-@o-rg!",
		"jay-@o.rg!",
		"jay-@o$rg!",
		".@org!",
		"af.@org!",
		"afa@org.!",
		"@@org!",
		"@org!",
		"ab@de!",
		"abc@d!",
		"jay",
		"yjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayx!",
		"yjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjay1@yjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjay1!",
		"yjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjay1@yjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjay!",
		"yjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjay@yjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjayjay1!",
		".jay!",
		"-jay!",
		"jay-!",
		"jay.!",
		"jay@-org!",
		"jay@.org!",
		"jay@org-!",
		"jay@org.!",
	}

	for _, address := range invalidAddresses {
		assert.False(t, IsValidAddress(address), address)
	}
}

func Test_Hashes(t *testing.T) {
	a, err := NewAddress("joshua@bitmaelum!")
	assert.NoError(t, err)
	assert.Equal(t, "68433f537c388686507649f75395a90c2d3b267eb2dc21f2443ca9006d31ad39", a.Hash().String())
	assert.Equal(t, "fc52fabe94c0e037d2df4498e87481a6438960c9f73d517584a7a5c564535ac4", a.LocalHash())
	assert.Equal(t, "49aa67181f4a3176f9b65605390bb81126e8ff1f6d03b1bd220c53e7a6b36d3e", a.OrgHash())

	a, err = NewAddress("joshua!")
	assert.NoError(t, err)
	assert.Equal(t, "66c94b6643ada5661b2d940eb87502b5af0f47f40fd45ce0fa125502dfa9c1ee", a.Hash().String())
	assert.Equal(t, "fc52fabe94c0e037d2df4498e87481a6438960c9f73d517584a7a5c564535ac4", a.LocalHash())
	assert.Equal(t, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", a.OrgHash())
}


func Test_Remainders(t *testing.T) {
	a, err := NewAddress("john@foobar!")
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x6a, 0x6f, 0x68, 0x6e, 0x40, 0x66, 0x6f, 0x6f, 0x62, 0x61, 0x72, 0x21}, a.Bytes())

	a, err = NewAddress("john!")
	assert.NoError(t, err)
	assert.False(t, a.HasOrganisationPart())
	assert.Equal(t, "john!", a.String())

	a, err = NewAddress("john@acme!")
	assert.NoError(t, err)
	assert.True(t, a.HasOrganisationPart())
}
