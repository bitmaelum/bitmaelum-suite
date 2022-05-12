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

package pow

import (
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

var (
	errChallengeNotFound = errors.New("challenge not found")
	errChallengeExpired  = errors.New("challenge expired")
)

type boltStorage struct {
	client *bolt.DB
}

// BucketName is the bucket name to store the invitations on the bolt db
const BucketName = "pow"

// BoltDBFile is the filename to store the boltdb database
const BoltDBFile = "pow.db"

// NewBolt initializes a new repository
func NewBolt(dbpath *string) Storable {
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

func (b *boltStorage) Retrieve(challenge string) (*ProofOfWork, error) {
	pow := &ProofOfWork{}

	err := b.client.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			logrus.Trace("challenge not found in BOLT: ", challenge, nil)
			return errChallengeNotFound
		}

		data := bucket.Get([]byte(challenge))
		if data == nil {
			logrus.Trace("challenge not found in BOLT: ", data, nil)
			return errChallengeNotFound
		}

		err := json.Unmarshal([]byte(data), &pow)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if pow.Expires.Unix() < internal.TimeNow().Unix() {
		return nil, errChallengeExpired
	}

	return pow, nil
}

func (b *boltStorage) Store(pow *ProofOfWork) error {
	data, err := json.Marshal(pow)
	if err != nil {
		return err
	}

	err = b.client.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			logrus.Trace("unable to create bucket on BOLT: ", BucketName, err)
			return err
		}

		return bucket.Put([]byte(pow.Challenge), data)
	})

	return err
}

func (b *boltStorage) Remove(challenge string) error {
	return b.client.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			logrus.Trace("unable to delete challenge, challenge not found in BOLT: ", challenge)
			return nil
		}

		return bucket.Delete([]byte(challenge))
	})
}
