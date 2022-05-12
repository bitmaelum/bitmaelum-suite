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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDurationCorrect(t *testing.T) {
	dataProvider := []struct {
		s  string
		ms int64
	}{
		{s: "1d", ms: 24 * 3600 * 1000},
		{s: "10d", ms: 240 * 3600 * 1000},
		{s: "10d1h", ms: 241 * 3600 * 1000},
		{s: "10d1h1m", ms: (10*24*3600 + 1*3600 + 60) * 1000},
		{s: "1d1y", ms: 31622400000},
		{s: "1d1m", ms: 86460000},
		{s: "1m1d", ms: 86460000},
		{s: "5w141m", ms: 3032460000},
	}

	for _, entry := range dataProvider {
		d, err := ParseDuration(entry.s)
		assert.NoError(t, err)
		assert.Equal(t, entry.ms, d.Milliseconds())
	}
}

func TestDurationIncorrect(t *testing.T) {
	dataProvider := []string{
		"1h 1h",
		"",
		"1A",
		"-151w",
		"0w",
		"foobar",
		"5P1G",
		"1d1w4d",
	}

	for _, entry := range dataProvider {
		d, err := ParseDuration(entry)
		assert.Error(t, err)
		assert.Equal(t, d, time.Duration(0))
	}
}
