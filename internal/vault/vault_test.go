package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	v, err := New("", []byte{})
	assert.NoError(t, err)
	assert.NotNil(t, v)
}
