package address

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidAddress(t *testing.T) {
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

func TestInvalidAddress(t *testing.T) {
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

func TestHashes(t *testing.T) {
	a, err := NewAddress("joshua@bitmaelum!")
	assert.NoError(t, err)
	assert.Equal(t, "6b024a4e51c0c4a30c3750115c66be776253880bb4af0f313e3bf2236e808840", a.Hash().String())
	assert.Equal(t, "fc52fabe94c0e037d2df4498e87481a6438960c9f73d517584a7a5c564535ac4", a.LocalHash().String())
	assert.Equal(t, "49aa67181f4a3176f9b65605390bb81126e8ff1f6d03b1bd220c53e7a6b36d3e", a.OrgHash().String())

	a, err = NewAddress("joshua!")
	assert.NoError(t, err)
	assert.Equal(t, "66c94b6643ada5661b2d940eb87502b5af0f47f40fd45ce0fa125502dfa9c1ee", a.Hash().String())
	assert.Equal(t, "fc52fabe94c0e037d2df4498e87481a6438960c9f73d517584a7a5c564535ac4", a.LocalHash().String())
	assert.Equal(t, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", a.OrgHash().String())

	a, _ = NewAddress("example!")
	assert.Equal(t, "2244643da7475120bf84d744435d15ea297c36ca165ea0baaa69ec818d0e952f", a.Hash().String())

	a, _ = NewAddress("jaytaph!")
	assert.Equal(t, "88667a68d0976d6c9106d4a68b4097026f0daeaec1aeb8351b096637679cf350", a.Hash().String())

	a, _ = NewAddress("bitmaelum!")
	assert.Equal(t, "dcc59b2c83ea86f1fd4d82f54acebfb083283553ef798380b18a0b5e512e668b", a.Hash().String())

	a, _ = NewAddress("hello@bitmaelum!")
	assert.Equal(t, "f3828fb0917561b49b2229953110e65785228f5302973ee52208a76bffc26aee", a.Hash().String())

	a, _ = NewAddress("hello@example!")
	assert.Equal(t, "a5098c40c4b7e272403f94d752026f45faeab26b4d67804c887969461b032074", a.Hash().String())
}

func TestRemainders(t *testing.T) {
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
