package core

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersion(t *testing.T) {
	buf := new(bytes.Buffer)
	WriteVersionInfo("foobar", buf)
	s := buf.String()

	assert.Contains(t, s, "foobar version")
}
