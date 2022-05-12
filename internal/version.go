// Copyright (c) 2022 BitMaelum Authors
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
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/coreos/go-semver/semver"
)

var (
	// BuildDate as filled in during compilation
	BuildDate string
	// GitCommit sha as filled in during compilation
	GitCommit string
	// VersionTag is the v0.0.0 version to use
	VersionTag = "0.0.0"
)

// Version is a structure with the current version of the software
var Version = semver.New(strings.Replace(VersionTag, "v", "", 1))

// WriteVersionInfo writes a string with all version information
func WriteVersionInfo(name string, w io.Writer) {
	s := fmt.Sprintf("%s version %d.%d.%d\nBuilt: %s\nCommit: %s", name, Version.Major, Version.Minor, Version.Patch, BuildDate, GitCommit)
	_, _ = w.Write([]byte(s))
}

// VersionString returns a string with all version information
func VersionString(name string) string {
	var b bytes.Buffer
	WriteVersionInfo(name, &b)

	return strings.Replace(b.String(), "\n", " * ", -1)
}
