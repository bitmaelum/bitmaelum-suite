package internal

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/afero/mem"
	"github.com/stretchr/testify/assert"
)

func TestSetLogging(t *testing.T) {
	std := logrus.StandardLogger()

	SetLogging("trace", "stdout")
	assert.Equal(t, logrus.TraceLevel, std.Level)
	SetLogging("debug", "stdout")
	assert.Equal(t, logrus.DebugLevel, std.Level)
	SetLogging("info", "stdout")
	assert.Equal(t, logrus.InfoLevel, std.Level)
	SetLogging("warning", "stdout")
	assert.Equal(t, logrus.WarnLevel, std.Level)
	SetLogging("error", "stdout")
	assert.Equal(t, logrus.ErrorLevel, std.Level)

	// incorrect setting
	SetLogging("debug", "stdout")
	SetLogging("foobar", "stdout")
	assert.Equal(t, logrus.ErrorLevel, std.Level)
}

func TestOutput(t *testing.T) {
	fs = afero.NewMemMapFs()
	std := logrus.StandardLogger()

	// Test output
	SetLogging("debug", "stdout")
	assert.Equal(t, os.Stdout, std.Out)

	SetLogging("debug", "stderr")
	assert.Equal(t, os.Stderr, std.Out)

	SetLogging("debug", "/path/to/file.log")
	ok, _ := afero.Exists(fs, "/path/to/file.log")
	assert.True(t, ok)
	assert.IsType(t, &mem.File{}, std.Out)
}

func TestSyslog(t *testing.T) {
	std := logrus.StandardLogger()

	SetLogging("debug", "syslog")
	assert.Equal(t, os.Stderr, std.Out)

	SetLogging("debug", "syslog:syslog.io:12345")
	assert.Equal(t, os.Stderr, std.Out)
}
