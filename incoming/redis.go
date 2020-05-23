package incoming

import (
    "fmt"
    "github.com/go-redis/redis/v8"
    "time"
)

type repo struct {
    client *redis.Client
}

func NewRedis(opts *redis.Options) IncomingRepository {
    return &repo{
        client: redis.NewClient(opts)
    }
}

func (r repo) Has(path string) (bool, error) {
    _, err := r.client.Get(r.client.Context(), path).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (r repo) Get(path string) (string, error) {
    return r.client.Get(r.client.Context(), path).Result()
}

func (r repo) Set(path string, value string, expiry time.Duration) error {
    return r.client.Set(r.client.Context(), path, value, expiry).Err()
}

