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

package imapgw

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

var (
	errNotFound = errors.New("not found")
)

type boltStorage struct {
	client *bolt.DB
}

// BucketName is the bucket name to store the flags on the bolt db
const BucketName = "flags"

// BoltDBFile is the filename to store the boltdb database
const BoltDBFile = "flags.db"

// NewBolt initializes a new repository
func NewBolt(dbpath *string) Storable {
	if _, err := os.Stat(*dbpath); os.IsNotExist(err) {
		os.Mkdir(*dbpath, os.ModeDir)
	}

	dbFile := filepath.Join(*dbpath, BoltDBFile)
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		logrus.Error("Unable to open filepath ", dbFile, err)
		return nil
	}

	return &boltStorage{
		client: db,
	}
}

func (b *boltStorage) Retrieve(messageid string) (flags []string) {
	b.client.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			logrus.Trace("messageid not found in BOLT: ", messageid)
			return errNotFound
		}

		data := bucket.Get([]byte(messageid))
		if data == nil {
			logrus.Trace("messageid not found in BOLT: ", messageid)
			return errNotFound
		}

		err := json.Unmarshal([]byte(data), &flags)
		if err != nil {
			return err
		}

		return nil
	})

	return
}

func (b *boltStorage) Store(messageid string, flags []string) error {
	data, err := json.Marshal(flags)
	if err != nil {
		return err
	}

	err = b.client.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			logrus.Trace("unable to create bucket on BOLT: ", BucketName, err)
			return err
		}

		return bucket.Put([]byte(messageid), data)
	})

	return err
}
