package storage

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisStorage struct {
	client *redis.Client
}

// NewRedis create a new redis storage repository
func NewRedis(opts *redis.Options) Storable {
	return &redisStorage{
		client: redis.NewClient(opts),
	}
}

func (r *redisStorage) Retrieve(challenge string) (*ProofOfWork, error) {
	data, err := r.client.Get(r.client.Context(), challenge).Result()
	if err != nil {
		return nil, err
	}

	pow := &ProofOfWork{}
	err = json.Unmarshal([]byte(data), pow)
	if err != nil {
		return nil, err
	}

	return pow, nil
}

func (r *redisStorage) Store(pow *ProofOfWork) error {
	return r.client.Set(r.client.Context(), pow.Challenge, pow, pow.Expires.Sub(time.Now())).Err()
}

func (r *redisStorage) Remove(challenge string) error {
	return r.client.Del(r.client.Context(), challenge).Err()
}
