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

package webhook

import (
	"encoding/json"

	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

type boltRepo struct {
	client     *bolt.DB
	BucketName string
}

// BoltDBFile is the filename to store the boltdb database
const BoltDBFile = "webhooks.db"

func (b boltRepo) FetchByHash(h hash.Hash) ([]Type, error) {
	var webhooks []Type

	err := b.client.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.BucketName))
		if bucket == nil {
			logrus.Trace("webhooks for account not found in BOLT: ", h, nil)
			return errWebhookNotFound
		}

		c := bucket.Cursor()
		for k, data := c.First(); k != nil; k, data = c.Next() {
			w := Type{}
			err := json.Unmarshal(data, &w)
			if err != nil {
				continue
			}

			if w.Account.String() == h.String() {
				webhooks = append(webhooks, w)
			}
		}

		return nil
	})

	if err != nil {
		return []Type{}, err
	}

	return webhooks, nil
}

// Fetch a key from the repository, or err
func (b boltRepo) Fetch(ID string) (*Type, error) {
	w := &Type{}

	err := b.client.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.BucketName))
		if bucket == nil {
			logrus.Trace("webhook not found in storage")
			return errWebhookNotFound
		}

		data := bucket.Get([]byte(ID))
		if data == nil {
			logrus.Trace("webhook not found in storage")
			return errWebhookNotFound
		}

		err := json.Unmarshal([]byte(data), w)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return w, nil
}

// Store the given key in the repository
func (b boltRepo) Store(w Type) error {
	data, err := json.Marshal(w)
	if err != nil {
		return err
	}

	err = b.client.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(b.BucketName))
		if err != nil {
			logrus.Trace("unable to create bucket on BOLT: ", b.BucketName, err)
			return err
		}

		return bucket.Put([]byte(w.ID), data)
	})

	return err
}

// Remove the given key from the repository
func (b boltRepo) Remove(w Type) error {
	return b.client.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.BucketName))
		if bucket == nil {
			logrus.Trace("unable to delete apikey, apikey not found in BOLT: ", w.ID, nil)
			return nil
		}

		return bucket.Delete([]byte(w.ID))
	})
}
