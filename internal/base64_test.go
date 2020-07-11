package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_EncodeDecde(t *testing.T) {
	b := Encode([]byte("foobar"))
	assert.Equal(t, []byte{0x5a, 0x6d, 0x39, 0x76, 0x59, 0x6d, 0x46, 0x79}, b)

	s, err := Decode(b)
	assert.NoError(t, err)
	assert.Equal(t, "foobar", string(s))
}
