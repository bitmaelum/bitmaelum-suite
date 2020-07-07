package core

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Checksum(t *testing.T) {
	var c message.Checksum

	c = Sha1([]byte("foobar"))
	assert.Equal(t, "sha1", c.Hash)
	assert.Equal(t, "8843d7f92416211de9ebb963ff4ce28125932878", c.Value)

	c = Sha256([]byte("foobar"))
	assert.Equal(t, "sha256", c.Hash)
	assert.Equal(t, "c3ab8ff13720e8ad9047dd39466b3c8974e592c2fa383d4a3960714caef0c4f2", c.Value)

	c = Md5([]byte("foobar"))
	assert.Equal(t, "md5", c.Hash)
	assert.Equal(t, "3858f62230ac3c915f300c664312c63f", c.Value)

	c = Crc32([]byte("foobar"))
	assert.Equal(t, "crc32", c.Hash)
	assert.Equal(t, "9ef61f95", c.Value)
}
