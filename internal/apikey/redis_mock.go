package apikey

import (
	"context"
	"time"
)

// RedisClientMock is a structure that can return mocked results from a redis. These results are queued by the Queue
// function, and fetched within the actual redis-methods. This system does not work directly when mocking redis
// but can be used by the redis-bridge.
type RedisClientMock struct {
	queue map[string][][]interface{}
}

// Queue will queue a set of return values that are returned when a specific method is called.
func (r *RedisClientMock) Queue(f string, args ...interface{}) {
	if r.queue == nil {
		r.queue = make(map[string][][]interface{})
	}

	r.queue[f] = append(r.queue[f], args)
}

// fetchFromQueue is called by the actual redis mock calls to fetch the correct
func (r *RedisClientMock) fetchFromQueue(f string) []interface{} {
	ret, ok := r.queue[f]
	if !ok {
		panic(f)
	}

	if len(r.queue[f]) > 0 {
		r.queue[f] = r.queue[f][1:]
	}

	return ret[0]
}


func (r *RedisClientMock) Del(ctx context.Context, keys ...string) (int64, error) {
	val := r.fetchFromQueue("del")
	return val[0].(int64), getError(val[1])
}

func (r *RedisClientMock) Get(ctx context.Context, key string) (string, error) {
	val := r.fetchFromQueue("get")
	return val[0].(string), getError(val[1])
}

func (r *RedisClientMock) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	val := r.fetchFromQueue("set")
	return val[0].(string), getError(val[1])
}

func (r *RedisClientMock) SMembers(ctx context.Context, key string) ([]string, error) {
	val := r.fetchFromQueue("smembers")
	return val[0].([]string), getError(val[1])
}

func (r *RedisClientMock) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	val := r.fetchFromQueue("sadd")
	return val[0].(int64), getError(val[1])
}

func (r *RedisClientMock) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	val := r.fetchFromQueue("srem")
	return val[0].(int64), getError(val[1])
}

// getError will return the value cast as error, or explicitly nil, because we can't cast nil to an error
func getError(v interface{}) error {
	if v == nil {
		return nil
	}

	return v.(error)
}
