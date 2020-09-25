package subscription

import (
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

func (r redisRepo) Has(sub *Subscription) bool {
	i, err := r.client.Exists(r.client.Context(), createKey(sub)).Result()

	return err == nil && i > 0
}

func (r redisRepo) Store(sub *Subscription) error {
	_, err := r.client.Set(r.client.Context(), createKey(sub), sub, 0).Result()

	return err
}

func (r redisRepo) Remove(sub *Subscription) error {
	_, err := r.client.Del(r.client.Context(), createKey(sub)).Result()

	return err
}
