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

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/go-redis/redis/v7"
)

type redisRepo struct {
	client  internal.RedisResultWrapper
	context context.Context
}

// NewRedisRepository initializes a new repository
func NewRedisRepository(opts *redis.Options) Repository {
	c := redis.NewClient(opts)

	return &redisRepo{
		client:  &internal.RedisBridge{Client: *c},
		context: c.Context(),
	}
}

func (r redisRepo) Has(sub *Subscription) bool {
	i, err := r.client.Exists(r.context, createKey(sub))

	return err == nil && i > 0
}

func (r redisRepo) Store(sub *Subscription) error {
	_, err := r.client.Set(r.context, createKey(sub), sub, 0)

	return err
}

func (r redisRepo) Remove(sub *Subscription) error {
	_, err := r.client.Del(r.context, createKey(sub))

	return err
}
