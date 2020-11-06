// Copyright (c) 2020 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package proof_of_work

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
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
