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

package ticket

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

//BucketName is the bucket name to store the tickets on the bolt db
const BucketName = "tickets"

//BoltDBFile is the filename to store the boltdb database
const BoltDBFile = "tickets.db"

// NewBoltRepository initializes a new repository
func NewBoltRepository(dbpath string) Repository {
	dbFile := filepath.Join(dbpath, BoltDBFile)
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		logrus.Error("Unable to open filepath ", dbFile, err)
		return nil
	}

	return &boltRepo{
		client: db,
	}
}

var errTicketNotFound = errors.New("ticket not found")

// Fetch a ticket from the repository, or err
func (b boltRepo) Fetch(ticketID string) (*Ticket, error) {
	logrus.Trace("Trying to fetch ticket from BOLT: ", ticketID)

	ticket := &Ticket{}
	err := b.client.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			logrus.Trace("ticket not found in BOLT: ", ticketID, nil)
			return errTicketNotFound
		}

		data := bucket.Get([]byte(ticketID))
		if data == nil {
			logrus.Trace("ticket not found in BOLT: ", data, nil)
			return errTicketNotFound
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
