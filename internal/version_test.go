package internal

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	buf := new(bytes.Buffer)
	WriteVersionInfo("foobar", buf)
	s := buf.String()

	assert.Contains(t, s, "foobar version")
}

func TestVersionString(t *testing.T) {
	assert.Equal(t, "foo version 0.0.1 * Built:  * Commit: ", VersionString("foo"))
}
