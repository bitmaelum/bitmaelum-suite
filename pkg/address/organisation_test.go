package address

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func Test_Organisation(t *testing.T) {
	h, err := NewOrganisationHash("foobar")
	assert.NoError(t, err)
	assert.Equal(t, "c3ab8ff13720e8ad9047dd39466b3c8974e592c2fa383d4a3960714caef0c4f2", h.String())

	h, err = NewOrganisationHash("foobar!")
	assert.Error(t, err)
	assert.Nil(t, h)

	h, err = NewOrganisationHash("f")
	assert.Error(t, err)
	assert.Nil(t, h)

	h, err = NewOrganisationHash("++--")
	assert.Error(t, err)
	assert.Nil(t, h)
}

func TestAddress_OrganisationHash(t *testing.T) {
	a, err := NewAddress("john@foobar!")
	assert.NoError(t, err)
	assert.Equal(t, "91af59f08863ff86a40bcc78e846c3cc4697ec4d52d606d50d1f2237fcd18523", a.OrganisationHash().String())

	a, err = NewAddress("joshua@foobar!")
	assert.NoError(t, err)
	assert.Equal(t, "91af59f08863ff86a40bcc78e846c3cc4697ec4d52d606d50d1f2237fcd18523", a.OrganisationHash().String())


	a, err = NewAddress("john!")
	assert.NoError(t, err)
	assert.Equal(t, "3b7546ed79e3e5a7907381b093c5a182cbf364c5dd0443dfa956c8cca271cc33", a.OrganisationHash().String())

	a, err = NewAddress("joshua!")
	assert.NoError(t, err)
	assert.Equal(t, "3b7546ed79e3e5a7907381b093c5a182cbf364c5dd0443dfa956c8cca271cc33", a.OrganisationHash().String())
}

func Test_Verify(t *testing.T) {
	// joshua@bitmaelum!
	assert.True(t, VerifyHash(
		"6b024a4e51c0c4a30c3750115c66be776253880bb4af0f313e3bf2236e808840",
		"fc52fabe94c0e037d2df4498e87481a6438960c9f73d517584a7a5c564535ac4",
		"49aa67181f4a3176f9b65605390bb81126e8ff1f6d03b1bd220c53e7a6b36d3e",
	))

	assert.False(t, VerifyHash(
		"68433f537c388686507649f75395a90c2d3b267eb2dc21f2443ca9006d31ad39",
		"fc52fabe94c0e037d2df4498e87481a6438960c9f73d517584a7a5c564535ac4",
		"0000000000000006f9b65605390bb81126e8ff1f6d03b1bd220c53e7a6b36d3e",
	))

	// joshua!
	assert.True(t, VerifyHash(
		"66c94b6643ada5661b2d940eb87502b5af0f47f40fd45ce0fa125502dfa9c1ee",
		"fc52fabe94c0e037d2df4498e87481a6438960c9f73d517584a7a5c564535ac4",
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	))

	assert.False(t, VerifyHash(
		"66c94b6643ada5661b2d940eb87502b5af0f47f40fd45ce0fa125502dfa9c1ee",
		"fc52fabe94c0e037d2df4498e87481a6438960c9f73d517584a7a5c564535ac4",
		"00000000000c1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	))
}
