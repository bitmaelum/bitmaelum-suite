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

package processor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateRetry(t *testing.T) {
	var d time.Duration

	d = getNextRetryDuration(0)
	assert.Equal(t, time.Duration(60000000000), d)

	d = getNextRetryDuration(1)
	assert.Equal(t, time.Duration(60000000000), d)

	d = getNextRetryDuration(8)
	assert.Equal(t, time.Duration(300000000000), d)

	d = getNextRetryDuration(16)
	assert.Equal(t, time.Duration(300000000000), d)

	d = getNextRetryDuration(17)
	assert.Equal(t, time.Duration(300000000000), d)

	d = getNextRetryDuration(18)
	assert.Equal(t, time.Duration(1800000000000), d)

	d = getNextRetryDuration(24)
	assert.Equal(t, time.Duration(1800000000000), d)

	d = getNextRetryDuration(25)
	assert.Equal(t, time.Duration(1800000000000), d)

	d = getNextRetryDuration(26)
	assert.Equal(t, time.Duration(3600000000000), d)

	d = getNextRetryDuration(31)
	assert.Equal(t, time.Duration(3600000000000), d)
}
