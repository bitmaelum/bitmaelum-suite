package core

import "github.com/coreos/go-semver/semver"

const (
	VersionMajor 	int64 = 0
	VersionMinor 	int64 = 0
	VersionPatch 	int64 = 1
	BuildDate		string = ""
	GitCommit 		string = ""
	GitBranch	 	string = ""
	GitState	 	string = ""
	GitSummary	 	string = ""
)

var Version = semver.Version{
	Major: VersionMajor,
	Minor: VersionMinor,
	Patch: VersionPatch,
}
