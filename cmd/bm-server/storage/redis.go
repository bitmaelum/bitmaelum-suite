package storage

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisStorage struct {
	client *redis.Client
}

// NewRedis creates a new redis storage repository
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
	// pow.Expires is an absolute time, we need delta time for redis
	var expiry time.Duration = 0
	if pow.Expires.Unix() > 0 {
		expiry = time.Until(pow.Expires)
	}

	data, err := json.Marshal(pow)
	if err != nil {
		return err
	}

	return r.client.Set(r.client.Context(), pow.Challenge, data, expiry).Err()
}

func (r *redisStorage) Remove(challenge string) error {
	return r.client.Del(r.client.Context(), challenge).Err()
}
