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

package testing

import (
	"context"
	"time"
)

// RedisClientMock is a structure that can return mocked results from a redis. These results are queued by the Queue
// function, and fetched within the actual redis-methods. This system does not work directly when mocking redis
// but can be used by the redis-bridge.
type RedisClientMock struct {
	queue map[string][][]interface{}
}

// Queue will queue a set of return values that are returned when a specific method is called.
func (r *RedisClientMock) Queue(f string, args ...interface{}) {
	if r.queue == nil {
		r.queue = make(map[string][][]interface{})
	}

	r.queue[f] = append(r.queue[f], args)
}

// fetchFromQueue is called by the actual redis mock calls to fetch the correct
func (r *RedisClientMock) fetchFromQueue(f string) []interface{} {
	ret, ok := r.queue[f]
	if !ok {
		panic(f)
	}

	if len(r.queue[f]) > 0 {
		r.queue[f] = r.queue[f][1:]
	}

	return ret[0]
}

// Del deletes a key
func (r *RedisClientMock) Del(ctx context.Context, keys ...string) (int64, error) {
	val := r.fetchFromQueue("del")
	return val[0].(int64), getError(val[1])
}

// Get retrieves a key
func (r *RedisClientMock) Get(ctx context.Context, key string) (string, error) {
	val := r.fetchFromQueue("get")
	return val[0].(string), getError(val[1])
}

// Exists checks if a key exists
func (r *RedisClientMock) Exists(ctx context.Context, key string) (int64, error) {
	val := r.fetchFromQueue("exists")
	return val[0].(int64), getError(val[1])
}

// Set stores a key
func (r *RedisClientMock) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	val := r.fetchFromQueue("set")
	return val[0].(string), getError(val[1])
}

// SMembers returns a set
func (r *RedisClientMock) SMembers(ctx context.Context, key string) ([]string, error) {
	val := r.fetchFromQueue("smembers")
	return val[0].([]string), getError(val[1])
}

// SAdd adds to a set
func (r *RedisClientMock) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	val := r.fetchFromQueue("sadd")
	return val[0].(int64), getError(val[1])
}

// SRem removes from a set
func (r *RedisClientMock) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	val := r.fetchFromQueue("srem")
	return val[0].(int64), getError(val[1])
}

// getError will return the value cast as error, or explicitly nil, because we can't cast nil to an error
func getError(v interface{}) error {
	if v == nil {
		return nil
	}

	return v.(error)
}
