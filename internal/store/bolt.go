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

package store

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

var (
	errPathNotFound           = errors.New("store: path not found")
	errParentNotFound         = errors.New("store: parent not found")
	errCannotRemoveCollection = errors.New("store: cannot remove collection")
)

// boltEntryType is the structure that we save to boltdb
type boltEntryType struct {
	Path      hash.Hash   `json:"path"`
	Parent    *hash.Hash  `json:"parent"`
	Data      []byte      `json:"data"`
	Timestamp int64       `json:"timestamp"`
	Entries   []hash.Hash `json:"entries"`
	Signature []byte      `json:"signature"`
}

type boltRepo struct {
	Clients map[string]*bolt.DB
	path    string
}

const (
	// BoltDBFile is the filename to store the boltdb database
	BoltDBFile = "store.db"
	// BucketName is the bucket name to save the sote
	BucketName = "store"
)

// NewBoltRepository initializes a new repository
func NewBoltRepository(accountsPath string) Repository {
	return &boltRepo{
		Clients: make(map[string]*bolt.DB),
		path:    accountsPath,
	}
}

// OpenDB will try and open the store database
func (b boltRepo) OpenDb(account hash.Hash) error {
	// Open file
	p := filepath.Join(b.path, account.String()[:2], account.String()[2:], BoltDBFile)
	logrus.Trace("opening boltdb file: ", p)

	opts := bolt.DefaultOptions
	opts.Timeout = 5 * time.Second
	db, err := bolt.Open(p, 0600, opts)
	if err != nil {
		logrus.Trace("error while opening boltdb: ", err)
		return err
	}

	// Store in cache
	b.Clients[account.String()] = db

	rootHash := hash.New(account.String() + "/")

	// Check if root exists and make one if it doesn't
	if !b.HasEntry(account, rootHash) {
		entry := &EntryType{
			Path:      rootHash,
			Timestamp: internal.TimeNow().Unix(),
		}

		err := b.SetEntry(account, *entry)
		if err != nil {
			logrus.Trace("error while setting root entry")
			return err
		}
	}

	return nil
}

// CloseDb will close the store database - if openened
func (b boltRepo) CloseDb(account hash.Hash) error {
	// check if db exists
	db, ok := b.Clients[account.String()]
	if !ok {
		return nil
	}

	delete(b.Clients, account.String())
	return db.Close()
}

// HasEntry will return true when the database has the specific path present
func (b boltRepo) HasEntry(account, path hash.Hash) bool {
	client, err := b.getClientDb(account)
	if err != nil {
		return false
	}

	err = client.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			return errPathNotFound
		}

		data := bucket.Get(path.Byte())
		if data == nil {
			return errPathNotFound
		}

		return nil
	})

	return err == nil
}

// GetEntry will return the given entry
func (b boltRepo) GetEntry(account, path hash.Hash, recursive bool, since time.Time) (*EntryType, error) {
	client, err := b.getClientDb(account)
	if err != nil {
		return nil, err
	}

	entry := &EntryType{}

	err = client.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			return errPathNotFound
		}

		dbEntry := getFromBucket(bucket, path)
		if dbEntry == nil {
			return errPathNotFound
		}

		// Don't return this node when we are filtering on time
		if time.Unix(dbEntry.Timestamp, 0).Before(since) {
			return nil
		}

		entry = &EntryType{
			Path:        dbEntry.Path,
			Parent:      dbEntry.Parent,
			Data:        dbEntry.Data,
			Timestamp:   dbEntry.Timestamp,
			HasChildren: len(dbEntry.Entries) > 0,
			Signature:   dbEntry.Signature,
		}

		if !recursive {
			return nil
		}

		// iterate all child entries
		for _, childHash := range dbEntry.Entries {
			childEntry, err := b.GetEntry(account, childHash, recursive, since)
			if err != nil {
				continue
			}

			entry.Children = append(entry.Children, *childEntry)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (b boltRepo) SetEntry(account hash.Hash, entry EntryType) error {
	client, err := b.getClientDb(account)
	if err != nil {
		return err
	}

	// Check if parent exists
	if entry.Parent != nil && !b.HasEntry(account, *entry.Parent) {
		return errParentNotFound
	}

	// Update entry and tree back to root with this timestamp
	lastUpdateTimestamp := internal.TimeNow().Unix()

	// Update entry values
	entry.Timestamp = lastUpdateTimestamp

	return client.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			return err
		}

		// convert to bolt entry
		dbEntry := boltEntryType{
			Path:      entry.Path,
			Parent:    entry.Parent,
			Data:      entry.Data,
			Timestamp: entry.Timestamp,
			Signature: entry.Signature,
		}

		err = putInBucket(bucket, dbEntry)
		if err != nil {
			return err
		}

		// Update parent entry
		if entry.Parent != nil {
			parentDbEntry := getFromBucket(bucket, *entry.Parent)
			parentDbEntry.Entries = addToEntries(parentDbEntry.Entries, entry.Path)

			err := putInBucket(bucket, *parentDbEntry)
			if err != nil {
				return err
			}

			// Update all parents timestamps
			err = updateTreeTimestamps(bucket, dbEntry, lastUpdateTimestamp)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func getFromBucket(bucket *bolt.Bucket, path hash.Hash) *boltEntryType {
	data := bucket.Get(path.Byte())
	if data == nil {
		return nil
	}

	entry := &boltEntryType{}
	err := json.Unmarshal(data, &entry)
	if err != nil {
		return nil
	}

	return entry
}

func putInBucket(bucket *bolt.Bucket, entry boltEntryType) error {
	buf, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return bucket.Put(entry.Path.Byte(), buf)
}

// RemoveEntry will remove the path from the database, and update the collection tree
func (b boltRepo) RemoveEntry(account, path hash.Hash, recursive bool) error {
	client, err := b.getClientDb(account)
	if err != nil {
		return err
	}

	entry, err := b.GetEntry(account, path, false, time.Time{})
	if err != nil {
		return errPathNotFound
	}

	// @TODO: recursive deletion is not yet supported
	if entry.HasChildren {
		return errCannotRemoveCollection
	}

	return client.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			logrus.Trace("unable to create bucket on BOLT: ", BucketName, err)
			return err
		}

		dbEntry := getFromBucket(bucket, entry.Path)

		// Remove actual entry
		err = bucket.Delete(entry.Path.Byte())
		if err != nil {
			return err
		}

		// Update parent entry
		if entry.Parent != nil {
			parentDbEntry := getFromBucket(bucket, *entry.Parent)
			parentDbEntry.Entries = removeFromEntries(parentDbEntry.Entries, entry.Path)
			err := putInBucket(bucket, *parentDbEntry)
			if err != nil {
				return err
			}

			// Update all parents
			lastUpdateTimestamp := internal.TimeNow().Unix()
			return updateTreeTimestamps(bucket, *dbEntry, lastUpdateTimestamp)
		}

		return nil
	})
}

// getClientDB will open or create the account's store database
func (b boltRepo) getClientDb(account hash.Hash) (*bolt.DB, error) {
	// Fetch db file from cache
	db, ok := b.Clients[account.String()]
	if ok {
		return db, nil
	}

	// Open/create if not found in cache
	err := b.OpenDb(account)
	if err != nil {
		return nil, err
	}

	return b.Clients[account.String()], nil
}

// addToEntries will add the path, but only when it's not yet present in the list
func addToEntries(entries []hash.Hash, path hash.Hash) []hash.Hash {
	for i := range entries {
		if entries[i].String() == path.String() {
			return entries
		}
	}

	return append(entries, path)
}

// removeFromEntries will add the path, but only when it's not yet present in the list
func removeFromEntries(entries []hash.Hash, path hash.Hash) []hash.Hash {
	// Find element in list
	found := -1
	for i := range entries {
		if entries[i].String() == path.String() {
			found = i
		}
	}

	if found == -1 {
		return entries
	}

	return append(entries[:found], entries[found+1:]...)
}

func updateTreeTimestamps(bucket *bolt.Bucket, initialDbEntry boltEntryType, ts int64) error {
	dbEntry := &initialDbEntry

	for dbEntry.Parent != nil {
		// Get parent entry
		dbEntry = getFromBucket(bucket, *dbEntry.Parent)
		if dbEntry == nil {
			return errParentNotFound
		}

		// Update this parent entry
		dbEntry.Timestamp = ts

		// Save back
		err := putInBucket(bucket, *dbEntry)
		if err != nil {
			return err
		}
	}

	return nil
}
