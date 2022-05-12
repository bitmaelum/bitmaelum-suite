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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWebhook(t *testing.T) {
	cfg := ConfigHTTP{
		URL: "https://foo.bar/test",
	}

	a, err := NewWebhook(hash.New("example!"), EventLocalDelivery, TypeHTTP, cfg)
	assert.NoError(t, err)

	u, err := uuid.Parse(a.ID)
	assert.NoError(t, err)
	assert.Equal(t, "VERSION_4", u.Version().String())
	assert.Equal(t, "2e4551de804e27aacf20f9df5be3e8cd384ed64488b21ab079fb58e8c90068ab", a.Account.String())
	assert.Equal(t, TypeHTTP, a.Type)
	assert.Equal(t, "{\"URL\":\"https://foo.bar/test\"}", string(a.Config))
	assert.False(t, a.Enabled)
	assert.Equal(t, EventLocalDelivery, a.Event)
}
