package subscription

import (
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

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
			return errors.New("subscription not found")
		}

		data := bucket.Get([]byte(createKey(sub)))
		if data == nil {
			logrus.Trace("subscription not found in BOLT: ", data, nil)
			return errors.New("subscription not found")
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
