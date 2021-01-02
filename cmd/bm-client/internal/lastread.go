// Copyright (c) 2021 BitMaelum Authors
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
