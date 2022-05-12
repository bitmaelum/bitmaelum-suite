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

package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var called = 0

func barFunc() string {
	return "this is bar"
}

func TestContainer(t *testing.T) {
	c := NewContainer()

	// We set a function that we can call here
	c.SetShared("foo", func() (interface{}, error) {
		return func() string {
			return "this is foo"
		}, nil
	})

	svc, ok := c.Get("foo").(func() string)
	assert.True(t, ok)
	assert.Equal(t, "this is foo", svc())

	c.SetShared("foo", func() (interface{}, error) { called++; return barFunc, nil })
	_, _ = c.Get("foo").(func() string)
	assert.Equal(t, 1, called)
	_, _ = c.Get("foo").(func() string)
	assert.Equal(t, 1, called)
	_, _ = c.Get("foo").(func() string)
	assert.Equal(t, 1, called)
	_, _ = c.Get("foo").(func() string)
	assert.Equal(t, 1, called)

	c.SetNonShared("foo", func() (interface{}, error) { called++; return barFunc, nil })
	_, _ = c.Get("foo").(func() string)
	assert.Equal(t, 2, called)
	_, _ = c.Get("foo").(func() string)
	assert.Equal(t, 3, called)
	_, _ = c.Get("foo").(func() string)
	assert.Equal(t, 4, called)
	_, _ = c.Get("foo").(func() string)
	assert.Equal(t, 5, called)
}
