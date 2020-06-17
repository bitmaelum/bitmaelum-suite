package core

import (
	"fmt"
	"github.com/coreos/go-semver/semver"
	"io"
)

const (
	VersionMajor 	int64 = 0
	VersionMinor 	int64 = 0
	VersionPatch 	int64 = 1
)

var BuildDate, GitCommit string


var Version = semver.Version{
	Major: VersionMajor,
	Minor: VersionMinor,
	Patch: VersionPatch,
}

func WriteVersionInfo(name string, w io.Writer) {
	s := fmt.Sprintf("%s version %d.%d.%d\nBuilt: %s\nCommit: %s", name, VersionMajor, VersionMinor, VersionPatch, BuildDate, GitCommit)
	w.Write([]byte(s))
}
