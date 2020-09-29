package vault

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_New(t *testing.T) {
	v, err := New("", []byte{})
	assert.NoError(t, err)
	assert.NotNil(t, v)
}
