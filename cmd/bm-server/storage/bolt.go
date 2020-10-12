package storage

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

type boltStorage struct {
	client *bolt.DB
}

//BucketName is the bucket name to store the invitations on the bolt db
const BucketName = "pow"

//BoltDBFile is the filename to store the boltdb database
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
			return errors.New("challenge not found")
		}

		data := bucket.Get([]byte(challenge))
		if data == nil {
			logrus.Trace("challenge not found in BOLT: ", data, nil)
			return errors.New("challenge not found")
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

	if pow.Expires.Unix() < time.Now().Unix() {
		return nil, errors.New("expired")
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
