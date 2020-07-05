package core

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLogo(t *testing.T) {
	assert.NotContains(t, GetMonochromeASCIILogo(), "\033[")

	_ = os.Setenv("COLORTERM", "truecolor")
	assert.Contains(t, GetASCIILogo(), "\033[38;5;209m")

	_ = os.Setenv("COLORTERM", "d")
	_ = os.Setenv("TERM", "xterm")
	assert.Contains(t, GetASCIILogo(), "\033[32m")

	_ = os.Setenv("COLORTERM", "")
	_ = os.Setenv("TERM", "linux")
	assert.NotContains(t, GetASCIILogo(), "\033[")

}
