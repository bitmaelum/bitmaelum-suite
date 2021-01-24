package imap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSequenceSet(t *testing.T) {
	set := NewSequenceSet("1,2,3")
	assert.False(t, set.InSet(0))
	assert.True(t, set.InSet(1))
	assert.True(t, set.InSet(2))
	assert.True(t, set.InSet(3))
	assert.False(t, set.InSet(4))


	set = NewSequenceSet("1,2,4")
	assert.False(t, set.InSet(0))
	assert.True(t, set.InSet(1))
	assert.True(t, set.InSet(2))
	assert.False(t, set.InSet(3))
	assert.True(t, set.InSet(4))

	set = NewSequenceSet("1:4")
	assert.False(t, set.InSet(0))
	assert.True(t, set.InSet(1))
	assert.True(t, set.InSet(2))
	assert.True(t, set.InSet(3))
	assert.True(t, set.InSet(4))
	assert.False(t, set.InSet(5))

	set = NewSequenceSet("1:3,6:7,10:10,2,12,14,15")
	assert.False(t, set.InSet(0))
	assert.True(t, set.InSet(1))
	assert.True(t, set.InSet(2))
	assert.True(t, set.InSet(3))
	assert.False(t, set.InSet(4))
	assert.False(t, set.InSet(5))
	assert.True(t, set.InSet(6))
	assert.True(t, set.InSet(7))
	assert.False(t, set.InSet(8))
	assert.False(t, set.InSet(9))
	assert.True(t, set.InSet(10))
	assert.False(t, set.InSet(11))
	assert.True(t, set.InSet(12))
	assert.False(t, set.InSet(13))
	assert.True(t, set.InSet(14))
	assert.True(t, set.InSet(15))
	assert.False(t, set.InSet(16))

	set = NewSequenceSet("1:*")
	assert.False(t, set.InSet(0))
	assert.True(t, set.InSet(1))
	assert.True(t, set.InSet(2))
	assert.True(t, set.InSet(3))
	assert.True(t, set.InSet(5252))
	assert.True(t, set.InSet(25232))
	assert.True(t, set.InSet(1000000))
}
