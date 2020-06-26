package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ProofOfWork(t *testing.T) {
	pow := NewProofOfWork(8, []byte("john@example!"), 0)
	assert.Equal(t, 8, pow.Bits)
	assert.Equal(t, uint64(0), pow.Proof)
	assert.False(t, pow.HasDoneWork())
	assert.False(t, pow.Validate())

	pow.Work()
	assert.True(t, pow.HasDoneWork())
	assert.True(t, pow.Validate())
	assert.Equal(t, uint64(88), pow.Proof)

	pow = NewProofOfWork(8, []byte("jane@example!"), 171)
	assert.Equal(t, 8, pow.Bits)
	assert.Equal(t, uint64(171), pow.Proof)
	assert.True(t, pow.HasDoneWork())
	assert.True(t, pow.Validate())
}
