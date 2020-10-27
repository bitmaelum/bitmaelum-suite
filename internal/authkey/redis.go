// Copyright (c) 2020 BitMaelum Authors
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

package authkey

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/go-redis/redis/v8"
)

// We don't use redis clients directly, but a redis result wrapper which does the calling of .Result() for us. This
// makes testing and mocking redis clients possible. It serves no other purpose.

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

// FetchByHash will retrieve all keys for the given account
func (r redisRepo) FetchByHash(h string) ([]KeyType, error) {
	keys := []KeyType{}

	items, err := r.client.SMembers(r.context, h)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		key, err := r.Fetch(item)
		if err != nil {
			continue
		}

		keys = append(keys, *key)
	}

	return keys, nil
}

// Fetch a key from the repository, or err
func (r redisRepo) Fetch(fingerprint string) (*KeyType, error) {
	data, err := r.client.Get(r.context, createRedisKey(fingerprint))
	if data == "" || err != nil {
		return nil, errors.New("key not found")
	}

	key := &KeyType{}
	err = json.Unmarshal([]byte(data), &key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// Store the given key in the repository
func (r redisRepo) Store(apiKey KeyType) error {
	data, err := json.Marshal(apiKey)
	if err != nil {
		return err
	}

	_, err = r.client.SAdd(r.context, apiKey.AddrHash.String(), apiKey.Fingerprint)
	if err != nil {
		return err
	}

	_, err = r.client.Set(r.context, createRedisKey(apiKey.Fingerprint), data, 0)
	return err
}

// Remove the given key from the repository
func (r redisRepo) Remove(apiKey KeyType) error {
	_, err := r.client.SRem(r.context, apiKey.AddrHash.String(), apiKey.Fingerprint)
	if err != nil {
		return err
	}

	_, err = r.client.Del(r.context, createRedisKey(apiKey.Fingerprint))
	return err
}

// createRedisKey creates a key based on the given fingerprint. This is needed otherwise we might send
// any data as fingerprint to redis in order to extract other kind of data (and you don't want that).
func createRedisKey(fingerprint string) string {
	return fmt.Sprintf("auth-%s", fingerprint)
}
