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
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/xhit/go-str2duration/v2"
)

var errUnknownPermission = func(p string) error { return errors.New("unknown permission: " + p) }

// ValidDuration gets a time duration string and return the time duration. Accepts single int as days
func ValidDuration(ds string) (time.Duration, error) {
	if ds == "" {
		return 0, nil
	}

	vd, err := str2duration.ParseDuration(ds)
	if err != nil {
		days, err := strconv.Atoi(ds)
		if err != nil {
			return 0, err
		}

		vd = time.Duration(days*24) * time.Hour
	}

	return vd, nil
}

// CheckManagementPermissions checks all permission and returns an error when a permission is not valid
func CheckManagementPermissions(perms []string) error {
	for _, p := range perms {
		p = strings.ToLower(p)

		found := false
		for _, ap := range ManagementPermissions {
			if p == ap {
				found = true
			}
		}
		if !found {
			return errUnknownPermission(p)
		}
	}

	return nil
}

// CheckAccountPermissions checks all permission and returns an error when a permission is not valid
func CheckAccountPermissions(perms []string) error {
	for _, p := range perms {
		p = strings.ToLower(p)

		found := false
		for _, ap := range AccountPermissions {
			if p == ap {
				found = true
			}
		}
		if !found {
			return errUnknownPermission(p)
		}
	}

	return nil
}
