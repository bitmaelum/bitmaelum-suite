package internal

import (
	"bytes"
	"fmt"
	"github.com/coreos/go-semver/semver"
	"io"
	"strings"
)

const (
	versionMajor int64 = 0
	versionMinor int64 = 0
	versionPatch int64 = 1
)

var buildDate, gitCommit string

// // Version is a structure with the current version of the software
var Version = semver.Version{
	Major: versionMajor,
	Minor: versionMinor,
	Patch: versionPatch,
}

// WriteVersionInfo writes a string with all version information
func WriteVersionInfo(name string, w io.Writer) {
	s := fmt.Sprintf("%s version %d.%d.%d\nBuilt: %s\nCommit: %s", name, versionMajor, versionMinor, versionPatch, buildDate, gitCommit)
	_, _ = w.Write([]byte(s))
}

// VersionString returns a string with all version information
func VersionString(name string) string {
	var b bytes.Buffer
	WriteVersionInfo(name, &b)

	return strings.Replace(b.String(), "\n", " * ", -1)
}
