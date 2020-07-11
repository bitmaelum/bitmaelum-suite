package internal

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Checksum(t *testing.T) {
	r := bytes.NewBufferString("foobar")
	c, err := CalculateChecksums(r)
	assert.NoError(t, err)

	assert.Equal(t, "8843d7f92416211de9ebb963ff4ce28125932878", c["sha1"])
	assert.Equal(t, "c3ab8ff13720e8ad9047dd39466b3c8974e592c2fa383d4a3960714caef0c4f2", c["sha256"])
	assert.Equal(t, "3858f62230ac3c915f300c664312c63f", c["md5"])
	assert.Equal(t, "0a50261ebd1a390fed2bf326f2673c145582a6342d523204973d0219337f81616a8069b012587cf5635f6925f1b56c360230c19b273500ee013e030601bf2425", c["sha512"])
}
