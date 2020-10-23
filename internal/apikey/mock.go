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

package apikey

import (
	"encoding/json"
	"errors"
)

type mockRepo struct {
	keys map[string][]byte
	addr map[string]map[string]int
}

// NewMockRepository initializes a new repository
func NewMockRepository() Repository {
	return &mockRepo{
		keys: map[string][]byte{},
		addr: map[string]map[string]int{},
	}
}

// FetchByHash will retrieve all keys for the given account
func (r mockRepo) FetchByHash(h string) ([]KeyType, error) {
	var keys []KeyType

	items, ok := r.addr[h]
	if !ok {
		return nil, errors.New("not found")
	}

	for item := range items {
		key, err := r.Fetch(item)
		if err != nil {
			continue
		}

		keys = append(keys, *key)
	}

	return keys, nil
}

// Fetch a key from the repository, or err
func (r mockRepo) Fetch(ID string) (*KeyType, error) {
	data, ok := r.keys[ID]
	if data == nil || !ok {
		return nil, errors.New("key not found")
	}

	key := &KeyType{}
	err := json.Unmarshal(data, &key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// Store the given key in the repository
func (r mockRepo) Store(apiKey KeyType) error {
	data, err := json.Marshal(apiKey)
	if err != nil {
		return err
	}

	// Add to account set if an hash is given
	if apiKey.AddrHash != nil {
		h := apiKey.AddrHash.String()
		if r.addr[h] == nil {
			r.addr[h] = map[string]int{}
		}
		r.addr[h][apiKey.ID] = 1
	}

	r.keys[apiKey.ID] = data
	return err
}

// Remove the given key from the repository
func (r mockRepo) Remove(apiKey KeyType) error {
	if apiKey.AddrHash != nil {
		delete(r.addr[apiKey.ID], apiKey.AddrHash.String())
	}

	delete(r.keys, apiKey.ID)
	return nil
}
