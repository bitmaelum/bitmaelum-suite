package apikey

import (
	"context"
	"testing"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/elliotchance/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedis(t *testing.T) {
	repo := redisRepo{
		client:  redismock.NewMock(),
		context: context.Background(),
	}

	var (
		err error
		kt  *KeyType
		kts []KeyType
	)

	h1 := hash.Hash("set 1")
	h2 := hash.Hash("set 2")

	kt = &KeyType{
		ID:          "abc",
		ValidUntil:  time.Time{},
		Permissions: []string{"foobar"},
		Admin:       true,
		AddrHash:    &h1,
		Desc:        "test key",
	}

	// Store does calls to SAdd and Set, how do we mock this, and fake result back?
	err = repo.Store(*kt)
	assert.NoError(t, err)
	kt, err = repo.Fetch("abc")
	assert.NoError(t, err)
	kt, err = repo.Fetch("efg")
	assert.NoError(t, err)

	kt = &KeyType{
		ID:          "def",
		ValidUntil:  time.Time{},
		Permissions: []string{"foobar"},
		Admin:       true,
		AddrHash:    &h2,
		Desc:        "test key 1",
	}
	kt = &KeyType{
		ID:          "ghi",
		ValidUntil:  time.Time{},
		Permissions: []string{"foobar"},
		Admin:       true,
		AddrHash:    &h1,
		Desc:        "test key 2",
	}

	kts, err = repo.FetchByHash("set 1")
	assert.NoError(t, err)
	assert.NotNil(t, kts)
	kts, err = repo.FetchByHash("set 2")
	assert.NoError(t, err)
	kts, err = repo.FetchByHash("set not exists")
	assert.NoError(t, err)

	err = repo.Remove(*kt)
	assert.NoError(t, err)
	kts, err = repo.FetchByHash("set 1")
	assert.NoError(t, err)
}

func TestCreateRedisKey(t *testing.T) {
	assert.Equal(t, "apikey-foobar", createRedisKey("foobar"))
}
