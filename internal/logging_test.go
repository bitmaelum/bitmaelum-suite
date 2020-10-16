// Copyright (c) 2020 BitMaelum Authors
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
