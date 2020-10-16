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
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

type boltRepo struct {
	client *bolt.DB
}

const (
	apiKeyNotFound string = "apikey not found"
)

// BucketName is the bucket name to store the invitations on the bolt db
const BucketName = "apikeys"

// BoltDBFile is the filename to store the boltdb database
const BoltDBFile = "apikeys.db"

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

func (b boltRepo) FetchByHash(h string) ([]KeyType, error) {
	keys := []KeyType{}

	err := b.client.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			logrus.Trace("keys for account not found in BOLT: ", h, nil)
			return errors.New(apiKeyNotFound)
		}

		// @TODO: we iterate all keys, unmarshall them to see if we need to add on a list. Please refactor
		//  into something better.. :(
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)

			key := &KeyType{}
			err := json.Unmarshal(v, &key)
			if err != nil {
				continue
			}

			if key.AddrHash.String() == h {
				keys = append(keys, *key)
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
func (b boltRepo) Fetch(ID string) (*KeyType, error) {
	key := &KeyType{}

	err := b.client.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			logrus.Trace("apikey not found in BOLT: ", ID, nil)
			return errors.New(apiKeyNotFound)
		}

		data := bucket.Get([]byte(ID))
		if data == nil {
			logrus.Trace("apikey not found in BOLT: ", data, nil)
			return errors.New(apiKeyNotFound)
		}

		err := json.Unmarshal([]byte(data), &key)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return key, nil
}

// Store the given key in the repository
func (b boltRepo) Store(apiKey KeyType) error {
	data, err := json.Marshal(apiKey)
	if err != nil {
		return err
	}

	err = b.client.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			logrus.Trace("unable to create bucket on BOLT: ", BucketName, err)
			return err
		}

		return bucket.Put([]byte(apiKey.ID), data)
	})

	return err
}

// Remove the given key from the repository
func (b boltRepo) Remove(apiKey KeyType) error {
	return b.client.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			logrus.Trace("unable to delete apikey, apikey not found in BOLT: ", apiKey.ID, nil)
			return nil
		}

		return bucket.Delete([]byte(apiKey.ID))
	})
}
