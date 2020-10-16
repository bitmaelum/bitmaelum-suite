package apikey

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/go-redis/redis/v8"
)

// We don't use redis clients directly, but a redis result wrapper which does the calling of .Result() for us. This
// makes testing and mocking redis clients possible. It serves no other purpose.

type redisRepo struct {
	client  internal.RedisResultWrapper
	context context.Context
}

// NewRedisRepository initializes a new repository
func NewRedisRepository(opts *redis.Options) Repository {
	c := redis.NewClient(opts)

	return &redisRepo{
		client:  &internal.RedisBridge{Client: *c},
		context: c.Context(),
	}
}

// FetchByHash will retrieve all keys for the given account
func (r redisRepo) FetchByHash(h string) ([]KeyType, error) {
	keys := []KeyType{}

	items, err := r.client.SMembers(r.context, h)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		key, err := r.Fetch(item)
		if err != nil {
			continue
		}

		keys = append(keys, *key)
	}

	return keys, nil
}

// Fetch a key from the repository, or err
func (r redisRepo) Fetch(ID string) (*KeyType, error) {
	data, err := r.client.Get(r.context, createRedisKey(ID))
	if data == "" || err != nil {
		return nil, errors.New("key not found")
	}

	key := &KeyType{}
	err = json.Unmarshal([]byte(data), &key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// Store the given key in the repository
func (r redisRepo) Store(apiKey KeyType) error {
	data, err := json.Marshal(apiKey)
	if err != nil {
		return err
	}

	// Add to account set if an hash is given
	if apiKey.AddrHash != nil {
		_, _ = r.client.SAdd(r.context, apiKey.AddrHash.String(), apiKey.ID)
	}

	_, err = r.client.Set(r.context, createRedisKey(apiKey.ID), data, 0)
	return err
}

// Remove the given key from the repository
func (r redisRepo) Remove(apiKey KeyType) error {
	if apiKey.AddrHash != nil {
		_, _ = r.client.SRem(r.context, apiKey.AddrHash.String(), apiKey.ID)
	}

	_, err := r.client.Del(r.context, createRedisKey(apiKey.ID))
	return err
}

// createRedisKey creates a key based on the given ID. This is needed otherwise we might send any data as api-id
// to redis in order to extract other kind of data (and you don't want that).
func createRedisKey(id string) string {
	return fmt.Sprintf("apikey-%s", id)
}
