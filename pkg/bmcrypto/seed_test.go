package bmcrypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateKeypairFromSeed(t *testing.T) {
	s, priv1, pub1, err := GenerateKeypairWithSeed()
	assert.NoError(t, err)

	priv2, pub2, err := GenerateKeypairFromSeed(s)
	assert.NoError(t, err)

	assert.Equal(t, priv1, priv2)
	assert.Equal(t, pub1, pub2)
}
