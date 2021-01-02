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

package key

import (
	"testing"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/stretchr/testify/assert"
)

func TestMockRepo(t *testing.T) {
	repo := NewAPIMockRepository()

	h := hash.New("example!")
	k := NewAPIAccountKey(h, []string{"foo", "bar"}, time.Time{}, "my desc")
	err := repo.Store(k)
	assert.NoError(t, err)

	k2, err := repo.Fetch(k.ID)
	assert.NoError(t, err)
	assert.Equal(t, []string{"foo", "bar"}, k2.Permissions)
	assert.Equal(t, "my desc", k2.Desc)

	h = hash.New("example!")
	k = NewAPIAccountKey(h, []string{"foo", "bar"}, time.Time{}, "my desc 2")
	err = repo.Store(k)
	assert.NoError(t, err)

	keys, err := repo.FetchByHash(h.String())
	assert.NoError(t, err)
	assert.Len(t, keys, 2)

	err = repo.Remove(k)
	assert.NoError(t, err)

	keys, err = repo.FetchByHash(h.String())
	assert.NoError(t, err)
	assert.Len(t, keys, 1)
}
