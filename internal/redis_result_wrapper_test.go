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
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"
)

func TestRedisBridge(t *testing.T) {
	rd, err := miniredis.Run()
	assert.NoError(t, err)
	defer rd.Close()

	client := redis.NewClient(&redis.Options{
		Addr: rd.Addr(),
	})
	bridge := RedisBridge{
		Client: *client,
	}

	var (
		i  int64
		s  string
		sl []string
	)

	ctx := context.TODO()

	i, err = bridge.Exists(ctx, "foo")
	assert.NoError(t, err)
	assert.Equal(t, i, int64(0))

	s, err = bridge.Set(ctx, "foo", "bar", time.Duration(1*time.Hour))
	assert.NoError(t, err)
	assert.Equal(t, s, "OK")

	i, err = bridge.Exists(ctx, "foo")
	assert.NoError(t, err)
	assert.Equal(t, i, int64(1))

	s, err = bridge.Get(ctx, "foo")
	assert.NoError(t, err)
	assert.Equal(t, s, "bar")

	i, err = bridge.Del(ctx, "foo", "bar")
	assert.NoError(t, err)
	assert.Equal(t, i, int64(1))

	i, err = bridge.SAdd(ctx, "foo", "bar", "baz")
	assert.NoError(t, err)
	assert.Equal(t, i, int64(2))

	sl, err = bridge.SMembers(ctx, "foo")
	assert.NoError(t, err)
	assert.Equal(t, sl, []string{"bar", "baz"})

	i, err = bridge.SRem(ctx, "foo", "baz")
	assert.NoError(t, err)
	assert.Equal(t, i, int64(1))

	sl, err = bridge.SMembers(ctx, "foo")
	assert.NoError(t, err)
	assert.Equal(t, sl, []string{"bar"})
}
