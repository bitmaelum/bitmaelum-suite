package subscription

import (
	"context"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/go-redis/redis/v8"
)

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

func (r redisRepo) Has(sub *Subscription) bool {
	i, err := r.client.Exists(r.context, createKey(sub))

	return err == nil && i > 0
}

func (r redisRepo) Store(sub *Subscription) error {
	_, err := r.client.Set(r.context, createKey(sub), sub, 0)

	return err
}

func (r redisRepo) Remove(sub *Subscription) error {
	_, err := r.client.Del(r.context, createKey(sub))

	return err
}
