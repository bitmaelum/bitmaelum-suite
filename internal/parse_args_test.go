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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPermissions(t *testing.T) {
	var err error

	err = CheckManagementPermissions([]string{})
	assert.NoError(t, err)

	err = CheckManagementPermissions([]string{"flush"})
	assert.NoError(t, err)

	err = CheckManagementPermissions([]string{"flush", "invite"})
	assert.NoError(t, err)

	err = CheckManagementPermissions([]string{"foo"})
	assert.Error(t, err)

	err = CheckManagementPermissions([]string{"flush", "foo"})
	assert.Error(t, err)

	err = CheckAccountPermissions([]string{"get-headers"})
	assert.NoError(t, err)

	err = CheckAccountPermissions([]string{"flush"})
	assert.Error(t, err)
}

func TestValidDuration(t *testing.T) {
	var d time.Duration
	var err error

	d, err = ValidDuration("")
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(0), d)

	_, err = ValidDuration("foobar")
	assert.Error(t, err)

	d, err = ValidDuration("1")
	assert.NoError(t, err)
	assert.Equal(t, 24*time.Hour, d)

	d, err = ValidDuration("1d")
	assert.NoError(t, err)
	assert.Equal(t, 24*time.Hour, d)

	d, err = ValidDuration("2h5m1s")
	assert.NoError(t, err)
	assert.Equal(t, 7501*time.Second, d)
}
