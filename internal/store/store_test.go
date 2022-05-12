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

package store

import (
	"testing"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/stretchr/testify/assert"
)

func TestNewEntry(t *testing.T) {
	internal.SetMockTime(func() time.Time {
		return time.Date(2010, 01, 01, 12, 34, 56, 0, time.UTC)
	})

	e := NewEntry([]byte("foobar"))
	assert.Equal(t, []byte("foobar"), e.Data)
	assert.Equal(t, int64(1262349296), e.Timestamp)
}

func TestMarshalling(t *testing.T) {
	internal.SetMockTime(func() time.Time {
		return time.Date(2010, 01, 01, 12, 34, 56, 0, time.UTC)
	})

	e := NewEntry([]byte("foobar"))
	b, err := e.MarshalBinary()
	assert.NoError(t, err)
	assert.Equal(t, "{\"path\":\"\",\"parent\":null,\"data\":\"Zm9vYmFy\",\"timestamp\":1262349296,\"has_children\":false,\"entries\":null,\"signature\":null}", string(b))

	e2 := &EntryType{}
	err = e2.UnmarshalBinary(b)
	assert.NoError(t, err)
	assert.Equal(t, []byte("foobar"), e2.Data)
	assert.Equal(t, int64(1262349296), e2.Timestamp)

}
