// Copyright (c) 2022 BitMaelum Authors
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

package key

import (
	"context"
	"testing"
	"time"

	testing2 "github.com/bitmaelum/bitmaelum-suite/internal/testing"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/stretchr/testify/assert"
)

func TestRedis(t *testing.T) {
	m := &testing2.RedisClientMock{}

	storage := redisRepo{
		client:  m,
		context: context.Background(),
	}

	repo := &APIKeyRepository{
		storageRepo: storage,
	}

	var (
		err error
		kt  *APIKeyType
		kts []APIKeyType
	)

	h1 := hash.Hash("set 1")
	kt = &APIKeyType{
		ID:          "abc",
		AddressHash: &h1,
		Expires:     time.Time{},
		Permissions: []string{"foobar"},
		Admin:       true,
		Desc:        "test key",
	}

	m.Queue("set", "ok", nil)
	m.Queue("sadd", int64(1), nil)
	err = repo.Store(*kt)
	assert.NoError(t, err)

	m.Queue("get", "{\"key\":\"abc\",\"valid_until\":\"0001-01-01T00:00:00Z\",\"permissions\":[\"foobar\"],\"admin\":true,\"address_hash\":\"set 1\",\"description\":\"test key\"}", nil)
	kt2, err := repo.Fetch("abc")
	assert.NoError(t, err)
	assert.Equal(t, "abc", kt2.ID)
	assert.Equal(t, "test key", kt2.Desc)

	m.Queue("get", "", errKeyNotFound)
	kt2, err = repo.Fetch("efg")
	assert.Error(t, err)
	assert.Nil(t, kt2)

	m.Queue("get", "notjson", nil)
	kt2, err = repo.Fetch("efg")
	assert.Error(t, err)
	assert.Nil(t, kt2)

	m.Queue("smembers", []string{"foo", "bar"}, nil)
	m.Queue("get", "{\"key\":\"abc\",\"valid_until\":\"0001-01-01T00:00:00Z\",\"permissions\":[\"foobar\"],\"admin\":true,\"address_hash\":\"set 1\",\"description\":\"test key\"}", nil)
	m.Queue("get", "{\"key\":\"def\",\"valid_until\":\"0001-01-01T00:00:00Z\",\"permissions\":[\"foobar\"],\"admin\":true,\"address_hash\":\"set 1\",\"description\":\"test key 2\"}", nil)
	kts, err = repo.FetchByHash("set 1")
	assert.NoError(t, err)
	assert.Len(t, kts, 2)
	assert.Equal(t, "abc", kts[0].ID)
	assert.Equal(t, "def", kts[1].ID)

	m.Queue("srem", int64(1), nil)
	m.Queue("del", int64(1), nil)
	err = repo.Remove(*kt)
	assert.NoError(t, err)
}

func TestCreateRedisKey(t *testing.T) {
	m := &testing2.RedisClientMock{}
	r := redisRepo{
		client:    m,
		context:   context.Background(),
		KeyPrefix: "apikey",
	}

	assert.Equal(t, "apikey-foobar", r.createRedisKey("foobar"))
}
