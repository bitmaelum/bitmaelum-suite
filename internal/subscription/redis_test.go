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

package subscription

import (
	"context"
	"errors"
	"testing"

	testing2 "github.com/bitmaelum/bitmaelum-suite/internal/testing"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepo(t *testing.T) {
	m := &testing2.RedisClientMock{}

	repo := &redisRepo{
		client:  m,
		context: context.Background(),
	}

	sub := New(hash.Hash("from"), hash.Hash("to"), "sub")

	m.Queue("exists", int64(1), nil)
	assert.True(t, repo.Has(&sub))

	m.Queue("exists", int64(0), nil)
	assert.False(t, repo.Has(&sub))

	m.Queue("exists", int64(33), nil)
	assert.True(t, repo.Has(&sub))

	m.Queue("exists", int64(0), errors.New("key not exist"))
	assert.False(t, repo.Has(&sub))

	m.Queue("set", "foobar", nil)
	assert.NoError(t, repo.Store(&sub))

	m.Queue("set", "foobar", errors.New("error"))
	assert.Error(t, repo.Store(&sub))

	m.Queue("del", int64(1), nil)
	assert.NoError(t, repo.Remove(&sub))

	m.Queue("del", int64(0), errors.New("error"))
	assert.Error(t, repo.Remove(&sub))
}
