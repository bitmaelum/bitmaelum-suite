package internal

import (
	"io/ioutil"
	"strconv"
	"time"

	"github.com/mitchellh/go-homedir"
)

const readTimeFile = "~/.bm-lastread"

// GetReadTime will return the last saved reading time or 0 when no time-file is found
func GetReadTime() time.Time {
	p, err := homedir.Expand(readTimeFile)
	if err != nil {
		return time.Time{}
	}

	b, err := ioutil.ReadFile(p)
	if err != nil {
		return time.Time{}
	}

	ts, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return time.Time{}
	}

	return time.Unix(ts, 0)
}

// SaveReadTime will save the read time to disk
func SaveReadTime(t time.Time) {
	p, err := homedir.Expand(readTimeFile)
	if err != nil {
		return
	}

	ts := strconv.FormatInt(t.Unix(), 10)
	_ = ioutil.WriteFile(p, []byte(ts), 0600)
}
