package message

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Checksum(t *testing.T) {
	r := bytes.NewBufferString("foobar")
	c, err := CalculateChecksums(r)
	assert.NoError(t, err)

	assert.Equal(t, "8843d7f92416211de9ebb963ff4ce28125932878", c["ripemd160"])
	assert.Equal(t, "c3ab8ff13720e8ad9047dd39466b3c8974e592c2fa383d4a3960714caef0c4f2", c["sha256"])
}
