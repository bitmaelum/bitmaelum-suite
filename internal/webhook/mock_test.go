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

package webhook

import (
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/stretchr/testify/assert"
)

var exampleHash = hash.New("example!")

func TestMockRepo(t *testing.T) {
	repo := NewMockRepository()

	cfg := ConfigHTTP{
		URL: "https://foo.bar/test",
	}

	w, err := NewWebhook(exampleHash, EventLocalDelivery, TypeHTTP, cfg)
	assert.NoError(t, err)

	err = repo.Store(*w)
	assert.NoError(t, err)

	w2, err := repo.Fetch(w.ID)
	assert.NoError(t, err)
	assert.False(t, w2.Enabled)
	assert.Equal(t, EventLocalDelivery, w2.Event)
	assert.Equal(t, TypeHTTP, w2.Type)
	assert.Equal(t, "2e4551de804e27aacf20f9df5be3e8cd384ed64488b21ab079fb58e8c90068ab", w2.Account.String())

	w, err = NewWebhook(exampleHash, EventLocalDelivery, TypeHTTP, cfg)
	assert.NoError(t, err)
	err = repo.Store(*w)
	assert.NoError(t, err)

	hooks, err := repo.FetchByHash(exampleHash)
	assert.NoError(t, err)
	assert.Len(t, hooks, 2)
	assert.NotEqual(t, hooks[0].ID, hooks[1].ID)
	assert.Equal(t, hooks[0].Account.String(), hooks[1].Account.String())

	err = repo.Remove(*w)
	assert.NoError(t, err)

	hooks, err = repo.FetchByHash(exampleHash)
	assert.NoError(t, err)
	assert.Len(t, hooks, 1)
}
