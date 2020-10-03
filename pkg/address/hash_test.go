package address

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHash(t *testing.T) {
	var h *Hash

	h = NewHash("joshua!")
	assert.Equal(t, "2f92571b5567b4557b94ac5701fc48e552ba9970d6dac89f7c2ebce92f1cd836", h.String())

	h = NewHash("somethingelse")
	assert.Equal(t, "102d0177f8ce6b0f5ada79780043ec90440529a53c98bd6043419e64c1d4274f", h.String())
}

func TestHashFromString(t *testing.T) {
	var (
		h   *Hash
		err error
	)

	h, err = HashFromString("12345678")
	assert.Error(t, err)
	assert.Nil(t, h)

	h, err = HashFromString("102d0177f8ce6b0f5ada79780043ec90440529a53c98bd6043419e64c1d4274f")
	assert.NoError(t, err)
	assert.Equal(t, "102d0177f8ce6b0f5ada79780043ec90440529a53c98bd6043419e64c1d4274f", h.String())
}
