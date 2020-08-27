package encrypt

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

// Mock reader so we deterministic signature
type dummyReader struct{}

func (d *dummyReader) Read(b []byte) (n int, err error) {
	for i := range b {
		b[i] = 1
	}
	return len(b), nil
}

func TestSign(t *testing.T) {
	randReader = &dummyReader{}

	data, _ := ioutil.ReadFile("./testdata/mykey.pem")
	privKey, _ := NewPrivKey(string(data))

	b, err := Sign(*privKey, []byte("b2d31086f098254d32314438a863e61e"))
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x4b, 0xf0, 0xf6, 0x3b, 0x80, 0x9d, 0xf1, 0x38, 0xa, 0x4f, 0xd, 0x24, 0x9c, 0x48, 0x4e, 0xdd, 0x10, 0xd5, 0x6a, 0x6d, 0x74, 0x1a, 0x22, 0x22, 0x86, 0xde, 0x8a, 0x17, 0xbe, 0x3, 0x66, 0x81, 0xba, 0xbe, 0x22, 0xca, 0x28, 0x26, 0x71, 0xaa, 0x49, 0x28, 0x99, 0x82, 0xb3, 0x28, 0x37, 0x36, 0x60, 0xc4, 0xf4, 0x1f, 0x14, 0x77, 0xef, 0xbb, 0x6e, 0xc3, 0x44, 0x45, 0x72, 0xf2, 0x83, 0x51, 0x42, 0x11, 0x81, 0xbe, 0x3, 0xe0, 0xcb, 0x36, 0x51, 0xd0, 0x4e, 0xf4, 0x63, 0x24, 0xaf, 0xd8, 0x2c, 0x74, 0x2, 0xc9, 0x5e, 0x2c, 0xb6, 0x2a, 0x3e, 0xb8, 0x72, 0x82, 0x13, 0xb0, 0x7a, 0xf6, 0x16, 0x29, 0x7e, 0x90, 0x3e, 0x42, 0x6a, 0xa2, 0xb5, 0xf7, 0xf2, 0x77, 0x92, 0xbf, 0x30, 0xa1, 0x1c, 0xe4, 0x9b, 0xe1, 0xf0, 0x45, 0x6e, 0x82, 0x5c, 0x46, 0x92, 0x4b, 0x58, 0xbf, 0x6d, 0x75}, b)

	res, err := Verify(*privKey, []byte("b2d31086f098254d32314438a863e61e"), b)
	assert.NoError(t, err)
	assert.True(t, res)

	b[0] ^= 0x80
	res, err = Verify(*privKey, []byte("b2d31086f098254d32314438a863e61e"), b)
	assert.NoError(t, err)
	assert.False(t, res)
}
