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
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var durationRegex = regexp.MustCompile("([0-9]+[ywdhm])")
var durationRegexMatch = regexp.MustCompile("^([0-9]+[ywdhm])+$")

var errInvalidFormat = errors.New("invalid format")

// ParseDuration will parse a string into a duration. This works pretty much the same way as the regular duration parser,
// except with larger times.
//
// It supports:
//     3y       // 3 years
//     5w       // 5 weeks
//     4d       // 4 days
//    14h       // 14 hours
//    52m       // 52 minutes
//
// or any combination of them. Any order is possible, but no two of the same units are allowed (1d4m3d for example).
func ParseDuration(s string) (time.Duration, error) {
	// Because we check submatches, we must make sure the "whole" string matches too.. otherwise we could match things
	// like  "-141w" or "1h 1h"
	if !durationRegexMatch.MatchString(s) {
		return time.Duration(0), errInvalidFormat
	}

	// Find all substrings (5m etc)
	matches := durationRegex.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return time.Duration(0), errInvalidFormat
	}

	d := 0
	found := ""
	for _, match := range matches {
		match := match[1]

		unit := string(match[len(match)-1])

		// check if we haven't seen it already
		if strings.Contains(found, unit) {
			return time.Duration(0), errInvalidFormat
		}
		found = found + unit

		// Check if quantity is valid
		qty, err := strconv.Atoi(match[:len(match)-1])
		if err != nil {
			return time.Duration(0), errInvalidFormat
		}
		// no 0 quantity
		if qty == 0 {
			return time.Duration(0), errInvalidFormat
		}

		switch unit {
		case "y": // year
			d += qty * 365 * 24 * 3600
		case "w": // week
			d += qty * 7 * 24 * 3600
		case "d": // day
			d += qty * 24 * 3600
		case "h": // hour
			d += qty * 3600
		case "m": // minute
			d += qty * 60
		}
	}

	return time.Duration(d * int(time.Second)), nil
}
