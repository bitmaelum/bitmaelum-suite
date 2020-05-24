package incoming

import (
    "github.com/go-redis/redis/v8"
    "time"
)

type redisRepo struct {
    client *redis.Client
}

func NewRedisRepository(opts *redis.Options) Repository {
    return &redisRepo{
        client: redis.NewClient(opts),
    }
}

func (r *redisRepo) Has(path string) (bool, error) {
    _, err := r.client.Get(r.client.Context(), path).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (r *redisRepo) Get(path string) ([]byte, error) {
    s, _ := r.client.Get(r.client.Context(), path).Result()
    return []byte(s), nil
}

func (r *redisRepo) Create(path string, value []byte, expiry time.Duration) error {
    return r.client.Set(r.client.Context(), path, value, expiry).Err()
}

func (r *redisRepo) Remove(path string) error {
    return r.client.Del(r.client.Context(), path).Err()
}
