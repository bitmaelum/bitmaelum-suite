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

package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

type BoltRepo struct {
	client *bolt.DB
}

type MessageInfo struct {
	BoxID       string
	UIDValidity int
	UID         int
	MessageID   string
	Flags       []string
}

type BoxInfo struct {
	BoxID       string
	UIDValidity int
	HighestUID  int
	Uids        []int
}

const BoltDBFile = "imap.db"

var boltClient *bolt.DB = nil

// NewBoltRepository initializes a new repository
func NewBolt(dbpath string) *BoltRepo {
	if boltClient == nil {
		var err error
		dbFile := filepath.Join(dbpath, BoltDBFile)
		boltClient, err = bolt.Open(dbFile, 0600, nil)
		if err != nil {
			logrus.Error("Unable to open filepath ", dbFile, err)
			return nil
		}
	}

	return &BoltRepo{
		client: boltClient,
	}
}

func (b BoltRepo) GetBoxInfo(bucket, boxID string) BoxInfo {
	info := BoxInfo{}

	err := b.client.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		if bucket == nil {
			return errors.New("not found")
		}

		data := bucket.Get([]byte(fmt.Sprintf("boxinfo-%s", boxID)))
		if data == nil {
			return errors.New("not found")
		}

		err := json.Unmarshal(data, &info)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		info = BoxInfo{
			BoxID:       boxID,
			UIDValidity: int(crc32.ChecksumIEEE([]byte(boxID))),
			HighestUID:  1000,
		}

		_ = b.StoreBoxInfo(bucket, info)
	}

	return info
}

func (b BoltRepo) StoreBoxInfo(bucket string, info BoxInfo) error {
	return b.client.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			logrus.Trace("unable to create bucket on BOLT: ", bucket, err)
			return err
		}

		buf, err := json.Marshal(info)
		if err != nil {
			return err
		}

		key := fmt.Sprintf("boxinfo-%s", info.BoxID)
		return bucket.Put([]byte(key), buf)
	})
}

func (b BoltRepo) Store(bucket, boxID string, uidValidity, UID int, messageID string) (*MessageInfo, error) {
	var info = MessageInfo{}

	err := b.client.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			logrus.Trace("unable to create bucket on BOLT: ", bucket, err)
			return err
		}

		info = MessageInfo{
			BoxID:       boxID,
			UIDValidity: uidValidity,
			UID:         UID,
			MessageID:   messageID,
			Flags:       []string{"\\Unseen"},
		}
		buf, err := json.Marshal(info)
		if err != nil {
			return err
		}

		// store key -> message data
		key := fmt.Sprintf("%s:%d:%d", boxID, uidValidity, UID)
		_ = bucket.Put([]byte(key), buf)

		// store messageid -> key
		_ = bucket.Put([]byte(messageID), []byte(key))

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (b BoltRepo) FetchByMessageID(bucket, messageID string) (*MessageInfo, error) {
	var key []byte

	err := b.client.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		if bucket == nil {
			return errors.New("not found")
		}

		key = bucket.Get([]byte(messageID))
		if key == nil {
			return errors.New("not found")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(key), ":")
	uid, _ := strconv.Atoi(parts[1])
	uidv, _ := strconv.Atoi(parts[2])
	return b.Fetch(bucket, parts[0], uid, uidv)
}

func (b BoltRepo) Fetch(bucket, boxID string, uidValidity, UID int) (*MessageInfo, error) {
	info := &MessageInfo{}
	err := b.client.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		if bucket == nil {
			return errors.New("not found")
		}

		key := fmt.Sprintf("%s:%d:%d", boxID, uidValidity, UID)
		data := bucket.Get([]byte(key))
		if data == nil {
			return errors.New("not found")
		}

		err := json.Unmarshal(data, &info)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return info, nil
}

func (b BoltRepo) Remove(bucket, boxID string, uidValidity, UID int) {
	_ = b.client.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		if bucket == nil {
			return nil
		}

		key := fmt.Sprintf("%s:%d:%d", boxID, uidValidity, UID)
		return bucket.Delete([]byte(key))
	})
}
