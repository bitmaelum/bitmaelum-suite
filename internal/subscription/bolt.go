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
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

var errSubscriptionNotFound = errors.New("subscription not found")

type boltRepo struct {
	client *bolt.DB
}

//BucketName is the bucket name to store the invitations on the bolt db
const BucketName = "subscriptions"

//BoltDBFile is the filename to store the boltdb database
const BoltDBFile = "subscriptions.db"

// NewBoltRepository initializes a new repository
func NewBoltRepository(dbpath string) Repository {
	dbFile := filepath.Join(dbpath, BoltDBFile)
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		logrus.Error("Unable to open filepath ", dbFile, err)
		return nil
	}

	return boltRepo{
		client: db,
	}
}

func (b boltRepo) Has(sub *Subscription) bool {

	err := b.client.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			logrus.Trace("subscription not found in BOLT: ", createKey(sub), nil)
			return errSubscriptionNotFound
		}

		data := bucket.Get([]byte(createKey(sub)))
		if data == nil {
			logrus.Trace("subscription not found in BOLT: ", data, nil)
			return errSubscriptionNotFound
		}

		return nil
	})

	return err == nil
}

func (b boltRepo) Store(sub *Subscription) error {
	data, err := json.Marshal(sub)
	if err != nil {
		return err
	}

	err = b.client.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			logrus.Trace("unable to create bucket on BOLT: ", BucketName, err)
			return err
		}

		return bucket.Put([]byte(createKey(sub)), data)
	})

	return err
}

func (b boltRepo) Remove(sub *Subscription) error {

	err := b.client.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			logrus.Trace("unable to delete subscription, subscription not found in BOLT: ", createKey(sub), nil)
			return nil
		}

		return bucket.Delete([]byte(createKey(sub)))
	})

	return err
}
