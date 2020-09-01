package proofofwork

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ProofOfWork(t *testing.T) {
	pow := New(8, "john@example!", 0)
	assert.Equal(t, 8, pow.Bits)
	assert.Equal(t, uint64(0), pow.Proof)
	assert.False(t, pow.HasDoneWork())
	assert.False(t, pow.IsValid())

	pow.Work()
	assert.True(t, pow.HasDoneWork())
	assert.True(t, pow.IsValid())
	assert.Equal(t, uint64(149), pow.Proof)

	pow = New(8, "jane@example!", 98)
	assert.Equal(t, 8, pow.Bits)
	assert.Equal(t, uint64(98), pow.Proof)
	assert.True(t, pow.HasDoneWork())
	assert.True(t, pow.IsValid())
}

func TestGenerateWorkData(t *testing.T) {
	w, e := GenerateWorkData()
	assert.NoError(t, e)
	assert.NotEmpty(t, w)
}
