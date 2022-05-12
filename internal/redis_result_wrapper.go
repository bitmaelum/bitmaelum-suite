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

package internal

import (
	"context"
	"time"

	"github.com/go-redis/redis/v7"
)

/* Defines a wrapper for redis results. Redis-go uses a system where you generate redis-commands which you need
 * to resolve manually. Since we cannot (easily) mock these result calls, we create a wrapper struct that does the
 * actual resolving for us. This means that now we can mock the wrapper now instead of mocking redis itself.
 *
 * Also, since we don't use all the redis methods, we only need to add the ones we are using by adding them to the
 * RedisResultWrapper interface.
 */

// RedisResultWrapper This is the our redis repository. It only contains the methods we really need.
type RedisResultWrapper interface {
	Del(ctx context.Context, keys ...string) (int64, error)
	Get(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, key string) (int64, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error)
	SMembers(ctx context.Context, key string) ([]string, error)
	SAdd(ctx context.Context, key string, members ...interface{}) (int64, error)
	SRem(ctx context.Context, key string, members ...interface{}) (int64, error)
}

// RedisBridge forms a bridge between a redis-client and redis implementors. This is needed to get rid of
// the ".Result()" calls, which makes mocking and testing difficult.
type RedisBridge struct {
	Client redis.Client
}

// Del removes a key
func (r RedisBridge) Del(ctx context.Context, keys ...string) (int64, error) {
	return r.Client.Del(keys...).Result()
}

// Get fetches a key
func (r RedisBridge) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(key).Result()
}

// Exists checks if a key exists
func (r RedisBridge) Exists(ctx context.Context, key string) (int64, error) {
	return r.Client.Exists(key).Result()
}

// Set stores a key
func (r RedisBridge) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	return r.Client.Set(key, value, expiration).Result()
}

// SMembers returns a set
func (r RedisBridge) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.Client.SMembers(key).Result()
}

// SAdd adds a key to a set
func (r RedisBridge) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return r.Client.SAdd(key, members...).Result()
}

// SRem removes a key from a set
func (r RedisBridge) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return r.Client.SRem(key, members...).Result()
}
