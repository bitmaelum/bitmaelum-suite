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
	"encoding/json"
	"reflect"
)

type mockRepo struct {
	keys map[string][]byte
	addr map[string]map[string]int
}

// FetchByHash will retrieve all keys for the given account
func (r mockRepo) FetchByHash(h string, v interface{}) (interface{}, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return nil, errNeedsPointerValue
	}

	items, ok := r.addr[h]
	if !ok {
		return nil, errKeyNotFound
	}

	var keys []interface{}

	ve := reflect.TypeOf(v).Elem()
	for item := range items {
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
func (r mockRepo) Fetch(ID string, v interface{}) error {
	data, ok := r.keys[ID]
	if data == nil || !ok {
		return errKeyNotFound
	}

	return json.Unmarshal(data, &v)
}

// Store the given key in the repository
func (r mockRepo) Store(v GenericKey) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// Add to account set if an hash is given
	if v.GetAddressHash() != nil {
		h := v.GetAddressHash().String()
		if r.addr[h] == nil {
			r.addr[h] = map[string]int{}
		}
		r.addr[h][v.GetID()] = 1
	}

	r.keys[v.GetID()] = data
	return err
}

// Remove the given key from the repository
func (r mockRepo) Remove(v GenericKey) error {
	if v.GetAddressHash() != nil {
		delete(r.addr[v.GetID()], v.GetAddressHash().String())
	}

	delete(r.keys, v.GetID())
	return nil
}
