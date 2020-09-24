package ticket

import (
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

type boltRepo struct {
	client *bolt.DB
}

//BucketName is the bucket name to store the tickets on the bolt db
const BucketName = "tickets"

//BoltDBFile is the filename to store the boltdb database
const BoltDBFile = "tickets.db"

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

// Fetch a ticket from the repository, or err
func (b boltRepo) Fetch(ticketID string) (*Ticket, error) {
	logrus.Trace("Trying to fetch ticket from BOLT: ", ticketID)

	ticket := &Ticket{}
	err := b.client.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			logrus.Trace("ticket not found in BOLT: ", ticketID, nil)
			return errors.New("ticket not found")
		}

		data := bucket.Get([]byte(ticketID))
		if data == nil {
			logrus.Trace("ticket not found in BOLT: ", data, nil)
			return errors.New("ticket not found")
		}

		err := json.Unmarshal([]byte(data), &ticket)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return ticket, nil
}

// Store the given ticket in the repository
func (b boltRepo) Store(ticket *Ticket) error {
	logrus.Trace("Storing ticket in BOLT: ", ticket)

	return b.client.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			logrus.Trace("unable to create bucket on BOLT: ", BucketName, err)
			return err
		}

		buf, err := json.Marshal(ticket)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(ticket.ID), buf)
	})
}

// Remove the given ticket from the repository
func (b boltRepo) Remove(ticketID string) {

	_ = b.client.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			logrus.Trace("unable to delete ticket, ticket not found in BOLT: ", ticketID, nil)
			return nil
		}

		return bucket.Delete([]byte(ticketID))
	})
}
