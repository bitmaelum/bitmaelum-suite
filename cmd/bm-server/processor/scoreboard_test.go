package processor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScoreboard(t *testing.T) {
	assert.False(t, IsInScoreboard(0, "1234"))

	AddToScoreboard(0, "1234")
	assert.True(t, IsInScoreboard(0, "1234"))
	assert.False(t, IsInScoreboard(1, "1234"))
	assert.False(t, IsInScoreboard(0, "5678"))

	RemoveFromScoreboard(0, "1234")
	assert.False(t, IsInScoreboard(0, "1234"))
}
