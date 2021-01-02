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

package webhook

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

type redisRepo struct {
	client    internal.RedisResultWrapper
	context   context.Context
	KeyPrefix string // redis key prefix
}

// FetchByHash will retrieve all hooks for the given account
func (r redisRepo) FetchByHash(h hash.Hash) ([]Type, error) {

	items, err := r.client.SMembers(r.context, r.createRedisKey(h.String()))
	if err != nil {
		return []Type{}, errWebhookNotFound
	}

	var webhooks []Type

	for _, item := range items {
		w, err := r.Fetch(item)
		if err != nil {
			continue
		}

		webhooks = append(webhooks, *w)
	}

	return webhooks, nil
}

// Fetch a key from the repository, or err
func (r redisRepo) Fetch(ID string) (*Type, error) {
	data, err := r.client.Get(r.context, r.createRedisKey(ID))
	if data == "" || err != nil {
		return nil, errWebhookNotFound
	}

	w := &Type{}
	err = json.Unmarshal([]byte(data), w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Store the given key in the repository
func (r redisRepo) Store(w Type) error {
	data, err := json.Marshal(w)
	if err != nil {
		return err
	}

	// Add to account set
	_, _ = r.client.SAdd(r.context, r.createRedisKey(w.Account.String()), w.ID)

	_, err = r.client.Set(r.context, r.createRedisKey(w.ID), data, 0)
	return err
}

// Remove the given key from the repository
func (r redisRepo) Remove(w Type) error {
	// Remove from account set
	_, _ = r.client.SRem(r.context, r.createRedisKey(w.Account.String()), w.ID)

	_, err := r.client.Del(r.context, r.createRedisKey(w.ID))
	return err
}

// createRedisKey creates a key based on the given ID. This is needed otherwise we might send any data as webhook-id
// to redis in order to extract other kind of data (and you don't want that).
func (r redisRepo) createRedisKey(id string) string {
	return fmt.Sprintf("%s-%s", r.KeyPrefix, id)
}
