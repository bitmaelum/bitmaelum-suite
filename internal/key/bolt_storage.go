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
	"encoding/json"

	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

type boltRepo struct {
	client     *bolt.DB
	BucketName string
}

// BoltDBFile is the filename to store the boltdb database
const BoltDBFile = "keys.db"

func (b boltRepo) FetchByHash(h string, v interface{}) (interface{}, error) {
	var keys []interface{}

	err := b.client.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.BucketName))
		if bucket == nil {
			logrus.Trace("keys for account not found in BOLT: ", h, nil)
			return errKeyNotFound
		}

		// @TODO: we iterate all keys, unmarshall them to see if we need to add on a list. Please refactor
		//  into something better.. :(
		c := bucket.Cursor()
		for k, data := c.First(); k != nil; k, data = c.Next() {
			err := json.Unmarshal(data, &v)
			if err != nil {
				continue
			}

			if v.(GenericKey).GetAddressHash().String() == h {
				keys = append(keys, v.(GenericKey))
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return keys, nil
}

// Fetch a key from the repository, or err
func (b boltRepo) Fetch(ID string, v interface{}) error {
	err := b.client.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.BucketName))
		if bucket == nil {
			logrus.Trace("api key not found in storage")
			return errKeyNotFound
		}

		data := bucket.Get([]byte(ID))
		if data == nil {
			logrus.Trace("api key not found in storage")
			return errKeyNotFound
		}

		err := json.Unmarshal([]byte(data), v)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// Store the given key in the repository
func (b boltRepo) Store(v GenericKey) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = b.client.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(b.BucketName))
		if err != nil {
			logrus.Trace("unable to create bucket on BOLT: ", b.BucketName, err)
			return err
		}

		return bucket.Put([]byte(v.GetID()), data)
	})

	return err
}

// Remove the given key from the repository
func (b boltRepo) Remove(v GenericKey) error {
	return b.client.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.BucketName))
		if bucket == nil {
			logrus.Trace("unable to delete apikey, apikey not found in BOLT: ", v.GetID(), nil)
			return nil
		}

		return bucket.Delete([]byte(v.GetID()))
	})
}
