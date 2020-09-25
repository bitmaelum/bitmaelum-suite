package invite

import (
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"

	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

type boltRepo struct {
	client *bolt.DB
}

type invitation struct {
	Token      string
	Expiration int64
}

//BucketName is the bucket name to store the invitations on the bolt db
const BucketName = "invitations"

//BoltDBFile is the filename to store the boltdb database
const BoltDBFile = "invitations.db"

// NewBoltRepository initializes a new repository
func NewBoltRepository(dbpath *string) Repository {
	dbFile := filepath.Join(*dbpath, BoltDBFile)
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		logrus.Error("Unable to open filepath ", dbFile, err)
		return nil
	}

	return &boltRepo{
		client: db,
	}
}

// Create generate a new invitation and stores this in bolt
func (b *boltRepo) Create(addr address.HashAddress, expiry time.Duration) (string, error) {
	buff := make([]byte, 32)
	_, err := rand.Read(buff)
	if err != nil {
		return "", err
	}
	token := base64.StdEncoding.EncodeToString(buff)

	err = b.client.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			logrus.Trace("unable to create bucket on BOLT: ", BucketName, err)
			return err
		}

		i := &invitation{}
		i.Token = token
		i.Expiration = time.Now().Add(expiry).Unix()

		buf, err := json.Marshal(i)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(createInviteKey(addr)), buf)
	})

	if err != nil {
		return "", err
	}

	return token, nil
}

// Get retrieves an invite from bolt
func (b *boltRepo) Get(addr address.HashAddress) (string, error) {

	i := &invitation{}

	err := b.client.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			logrus.Trace("invite not found in BOLT: ", createInviteKey(addr), nil)
			return errors.New("invite not found")
		}

		data := bucket.Get([]byte(createInviteKey(addr)))
		if data == nil {
			logrus.Trace("invite not found in BOLT: ", data, nil)
			return errors.New("invite not found")
		}

		err := json.Unmarshal([]byte(data), &i)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if i.Expiration < time.Now().Unix() {
		logrus.Trace("invite is expired: ", i.Expiration, nil)
		b.Remove(addr)
		return "", errors.New("invitation expired")
	}

	return i.Token, nil
}

// Remove deletes an invite from bolt
func (b *boltRepo) Remove(addr address.HashAddress) error {

	return b.client.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			logrus.Trace("unable to delete ticket, invite not found in BOLT: ", createInviteKey(addr), nil)
			return nil
		}

		return bucket.Delete([]byte(createInviteKey(addr)))
	})
}
