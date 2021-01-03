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

package key

import (
	"path/filepath"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

// AuthKeyRepository is a repository to fetch and store auth keys
type AuthKeyRepository struct {
	storageRepo StorageBackend
}

// NewAuthKeyRedisRepository initializes a new Redis repository
func NewAuthKeyRedisRepository(opts *redis.Options) AuthKeyRepo {
	c := redis.NewClient(opts)

	// Setup generic repository
	repo := &redisRepo{
		client:    &internal.RedisBridge{Client: *c},
		context:   c.Context(),
		KeyPrefix: "authkey",
	}

	return &AuthKeyRepository{
		storageRepo: repo,
	}
}

// NewAuthBoltRepository initializes a new BoltDb repository
func NewAuthBoltRepository(dbPath string) AuthKeyRepo {
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

	return &AuthKeyRepository{
		storageRepo: repo,
	}
}

// NewAuthMockRepository initializes a new mock repository
func NewAuthMockRepository() AuthKeyRepo {
	repo := &mockRepo{
		keys: map[string][]byte{},
		addr: map[string]map[string]int{},
	}

	return &AuthKeyRepository{
		storageRepo: repo,
	}
}

// FetchByHash will fetch all api keys for the given address hash
func (a AuthKeyRepository) FetchByHash(h string) ([]AuthKeyType, error) {
	v := &AuthKeyType{}
	l, err := a.storageRepo.FetchByHash(h, v)

	var ll []AuthKeyType
	for _, item := range l.([]interface{}) {
		if item == nil {
			continue
		}

		p := item.(*AuthKeyType)
		ll = append(ll, *p)
	}

	return ll, err
}

// Fetch will get a single API key
func (a AuthKeyRepository) Fetch(ID string) (*AuthKeyType, error) {
	v := &AuthKeyType{}

	err := a.storageRepo.Fetch(ID, &v)
	if err != nil {
		return nil, errKeyNotFound
	}

	return v, err
}

// Store will store a single api key
func (a AuthKeyRepository) Store(v AuthKeyType) error {
	return a.storageRepo.Store(v)
}

// Remove will remove a single api key
func (a AuthKeyRepository) Remove(v AuthKeyType) error {
	return a.storageRepo.Remove(v)
}
