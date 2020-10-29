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

package key

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/bitmaelum/bitmaelum-suite/internal"
)

type redisRepo struct {
	client    internal.RedisResultWrapper
	context   context.Context
	KeyPrefix string // redis key prefix
}

// FetchByHash will retrieve all keys for the given account
func (r redisRepo) FetchByHash(h string, v interface{}) (interface{}, error) {
	// v needs to be a pointer
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return nil, errNeedsPointerValue
	}

	items, err := r.client.SMembers(r.context, h)
	if err != nil {
		return nil, errKeyNotFound
	}

	var keys []interface{}

	ve := reflect.TypeOf(v).Elem()
	for _, item := range items {
		// Create a new item based on the structure of v
		newItem := reflect.New(ve).Interface()

		err := r.Fetch(item, &newItem)
		if err != nil {
			continue
		}

		keys = append(keys, newItem.(interface{}))
	}

	return keys, nil
}

// Fetch a key from the repository, or err
func (r redisRepo) Fetch(ID string, v interface{}) error {
	data, err := r.client.Get(r.context, r.createRedisKey(ID))
	if data == "" || err != nil {
		return errKeyNotFound
	}

	err = json.Unmarshal([]byte(data), v)
	if err != nil {
		return err
	}

	return nil
}

// Store the given key in the repository
func (r redisRepo) Store(k GenericKey) error {
	data, err := json.Marshal(k)
	if err != nil {
		return err
	}

	// Add to account set if an hash is given
	if k.GetAddressHash() != nil {
		_, _ = r.client.SAdd(r.context, r.createRedisKey(k.GetAddressHash().String()), k.GetID())
	}

	_, err = r.client.Set(r.context, r.createRedisKey(k.GetID()), data, 0)
	return err
}

// Remove the given key from the repository
func (r redisRepo) Remove(k GenericKey) error {
	if k.GetAddressHash() != nil {
		_, _ = r.client.SRem(r.context, r.createRedisKey(k.GetAddressHash().String()), k.GetID())
	}

	_, err := r.client.Del(r.context, r.createRedisKey(k.GetID()))
	return err
}

// createRedisKey creates a key based on the given ID. This is needed otherwise we might send any data as api-id
// to redis in order to extract other kind of data (and you don't want that).
func (r redisRepo) createRedisKey(id string) string {
	return fmt.Sprintf("%s-%s", r.KeyPrefix, id)
}
