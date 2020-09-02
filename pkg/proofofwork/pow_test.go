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

func TestString(t *testing.T) {
	pow := New(8, "john@example!", 149)
	assert.True(t, pow.IsValid())
	assert.Equal(t, "8$am9obkBleGFtcGxlIQ==$149", pow.String())

	pow, err := NewFromString("8$am9obkBleGFtcGxlIQ==$149")
	assert.NoError(t, err)
	assert.True(t, pow.IsValid())
	assert.Equal(t, uint64(149), pow.Proof)
	assert.Equal(t, 8, pow.Bits)
	assert.Equal(t, "john@example!", pow.Data)

	pow, err = NewFromString("8$am9obkBleGFtcGxlIQ==$12431241")
	assert.NoError(t, err)
	assert.False(t, pow.IsValid())
	assert.Equal(t, uint64(12431241), pow.Proof)
	assert.Equal(t, 8, pow.Bits)
	assert.Equal(t, "john@example!", pow.Data)

	pow, err = NewFromString("8$am9obkBleGFtcGxlI")
	assert.Error(t, err)
	assert.Nil(t, pow)

	pow, err = NewFromString("3$a$1")
	assert.NoError(t, err)
	assert.False(t, pow.IsValid())
	assert.Equal(t, uint64(1), pow.Proof)
	assert.Equal(t, 3, pow.Bits)
	assert.Equal(t, "", pow.Data)

	pow, err = NewFromString("8$a$b")
	assert.Error(t, err)
	assert.Nil(t, pow)
}
