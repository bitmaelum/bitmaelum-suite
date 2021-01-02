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

package webhook

import (
	"path/filepath"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

// Storage is the main interface for fetching and storing webhook data
type Storage interface {
	FetchByHash(h hash.Hash) ([]Type, error)
	Fetch(ID string) (*Type, error)
	Store(w Type) error
	Remove(w Type) error
}

// Repository is a repository to fetch and store webhooks
type Repository struct {
	storageRepo Storage
}

// NewRedisRepository initializes a new Redis repository
func NewRedisRepository(opts *redis.Options) *Repository {
	c := redis.NewClient(opts)

	// Setup generic repository
	repo := &redisRepo{
		client:    &internal.RedisBridge{Client: *c},
		context:   c.Context(),
		KeyPrefix: "webhook",
	}

	return &Repository{
		storageRepo: repo,
	}
}

// NewBoltRepository initializes a new BoltDb repository
func NewBoltRepository(dbPath string) *Repository {
	p := filepath.Join(dbPath, BoltDBFile)
	db, err := bolt.Open(p, 0600, nil)
	if err != nil {
		logrus.Error("Unable to open filepath ", p, err)
		return nil
	}

	repo := boltRepo{
		client:     db,
		BucketName: "auth",
	}

	return &Repository{
		storageRepo: repo,
	}
}

// NewMockRepository initializes a new mock repository
func NewMockRepository() *Repository {
	repo := &mockRepo{
		Webhooks: make(map[string]Type),
	}

	return &Repository{
		storageRepo: repo,
	}
}

// FetchByHash will fetch all api keys for the given address hash
func (r Repository) FetchByHash(h hash.Hash) ([]Type, error) {
	return r.storageRepo.FetchByHash(h)
}

// Fetch will get a single API key
func (r Repository) Fetch(ID string) (*Type, error) {
	return r.storageRepo.Fetch(ID)
}

// Store will store a single webhook
func (r Repository) Store(w Type) error {
	return r.storageRepo.Store(w)
}

// Remove will remove a single webhook
func (r Repository) Remove(w Type) error {
	return r.storageRepo.Remove(w)
}
