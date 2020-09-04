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

	assert.Equal(t, "a06e327ea7388c18e4740e350ed4e60f2e04fc41", c["ripemd160"])
	assert.Equal(t, "c3ab8ff13720e8ad9047dd39466b3c8974e592c2fa383d4a3960714caef0c4f2", c["sha256"])
}
