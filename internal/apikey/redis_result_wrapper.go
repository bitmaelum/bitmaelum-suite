package apikey

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)


/* Defines a wrapper for redis results. Redis-go uses a system where you generate redis-commands which you need
 * to resolve manually. Since we cannot (easily) mock these result calls, we create a wrapper struct that does the
 * actual resolving for us. This means that now we can mock the wrapper now instead of mocking redis itself.
 *
 * Also, since we don't use all the redis methods, we only need to add the ones we are using by adding them to the
 * RedisResultWrapper interface.
 */

// RedisResultWrapper This is the our redis repository. It only contains the methods we really need. THis
type RedisResultWrapper interface {
	Del(ctx context.Context, keys ...string) (int64, error)
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error)
	SMembers(ctx context.Context, key string) ([]string, error)
	SAdd(ctx context.Context, key string, members ...interface{}) (int64, error)
	SRem(ctx context.Context, key string, members ...interface{}) (int64, error)
}


type redisBridge struct {
	client redis.Client
}

func (r redisBridge) Del(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Del(ctx, keys...).Result()
}

func (r redisBridge) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r redisBridge) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	return r.client.Set(ctx, key, value, expiration).Result()
}

func (r redisBridge) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}

func (r redisBridge) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return r.client.SAdd(ctx, key, members...).Result()
}

func (r redisBridge) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return r.client.SRem(ctx, key, members...).Result()
}
