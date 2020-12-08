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

package webhook

import (
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

type mockRepo struct {
	Webhooks map[string]Type
}

// FetchByHash will retrieve all keys for the given account
func (r mockRepo) FetchByHash(h hash.Hash) ([]Type, error) {
	var webhooks []Type

	for _, w := range r.Webhooks {
		if w.Account.String() == h.String() {
			webhooks = append(webhooks, w)
		}
	}

	return webhooks, nil
}

// Fetch a key from the repository, or err
func (r mockRepo) Fetch(ID string) (*Type, error) {
	w, ok := r.Webhooks[ID]
	if !ok {
		return nil, errWebhookNotFound
	}

	return &w, nil
}

// Store the given key in the repository
func (r mockRepo) Store(w Type) error {
	r.Webhooks[w.ID] = w
	return nil
}

// Remove the given key from the repository
func (r mockRepo) Remove(w Type) error {
	delete(r.Webhooks, w.ID)
	return nil
}
