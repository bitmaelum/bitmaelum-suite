package apikey

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type redisRepo struct {
	client *redis.Client
}

// NewRedisRepository initializes a new repository
func NewRedisRepository(opts *redis.Options) Repository {
	return &redisRepo{
		client: redis.NewClient(opts),
	}
}

// Fetch a key from the repository, or err
func (r redisRepo) Fetch(ID string) (*KeyType, error) {
	data, err := r.client.Get(r.client.Context(), createRedisKey(ID)).Result()
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

	_, err = r.client.Set(r.client.Context(), createRedisKey(apiKey.ID), data, 0).Result()
	return err
}

// Remove the given key from the repository
func (r redisRepo) Remove(ID string) {
	_ = r.client.Del(r.client.Context(), createRedisKey(ID))
}

// createRedisKey creates a key based on the given ID. This is needed otherwise we might send any data as api-id
// to redis in order to extract other kind of data (and you don't want that).
func createRedisKey(id string) string {
	return fmt.Sprintf("apikey-%s", id)
}
