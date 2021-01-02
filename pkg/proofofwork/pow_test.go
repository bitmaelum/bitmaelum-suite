// Copyright (c) 2021 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package proofofwork

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProofOfWork(t *testing.T) {
	pow := New(8, "john@example!", 0)
	assert.Equal(t, 8, pow.Bits)
	assert.Equal(t, uint64(0), pow.Proof)
	assert.False(t, pow.HasDoneWork())
	assert.False(t, pow.IsValid())

	// Use a single core, otherwise we don't know which core will find the proof and thus what
	// the proof actually is.
	pow.Work(1)
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
	assert.Error(t, err)
	assert.Nil(t, pow)

	pow, err = NewFromString("8$a$b")
	assert.Error(t, err)
	assert.Nil(t, pow)

	pow, err = NewFromString("abc$am9obkBleGFtcGxlIQ==$149")
	assert.Error(t, err)
	assert.Nil(t, pow)

	pow, err = NewFromString("7$am9obkBleGFtcGxlIQ==$foobar")
	assert.Error(t, err)
	assert.Nil(t, pow)
}

func TestWorkData(t *testing.T) {
	randomReader = &dummyReader{}

	s, err := GenerateWorkData()
	assert.NoError(t, err)
	assert.Equal(t, "AQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQE=", s)
}

func TestMaxCores(t *testing.T) {
	assert.GreaterOrEqual(t, maxCores(), 1)

	pow := NewWithoutProof(10, "foobar")
	pow.WorkMulticore()
	assert.True(t, pow.IsValid())
}

func TestMarshalling(t *testing.T) {
	powBytes := []byte{0x22, 0x38, 0x24, 0x61, 0x6d, 0x39, 0x6f, 0x62, 0x6b, 0x42, 0x6c, 0x65, 0x47, 0x46, 0x74, 0x63, 0x47, 0x78, 0x6c, 0x49, 0x51, 0x3d, 0x3d, 0x24, 0x31, 0x34, 0x39, 0x22}

	pow, err := NewFromString("8$am9obkBleGFtcGxlIQ==$149")
	assert.NoError(t, err)

	b, err := pow.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, powBytes, b)

	pow = &ProofOfWork{}
	err = pow.UnmarshalJSON(powBytes)
	assert.NoError(t, err)
	assert.True(t, pow.IsValid())
	assert.Equal(t, uint64(149), pow.Proof)
	assert.Equal(t, 8, pow.Bits)
	assert.Equal(t, "john@example!", pow.Data)

	pow = &ProofOfWork{}
	err = pow.UnmarshalJSON([]byte("{111111}"))
	assert.Error(t, err)
	err = pow.UnmarshalJSON([]byte("\"fooobar\""))
	assert.Error(t, err)
}

func TestNewWithoutProof(t *testing.T) {
	pow := NewWithoutProof(22, "foo")
	assert.False(t, pow.IsValid())
}

type dummyReader struct{}

func (d *dummyReader) Read(b []byte) (n int, err error) {
	for i := range b {
		b[i] = 1
	}
	return len(b), nil
}
