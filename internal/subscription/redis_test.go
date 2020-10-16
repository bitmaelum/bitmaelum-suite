package subscription

import (
	"context"
	"errors"
	"testing"

	testing2 "github.com/bitmaelum/bitmaelum-suite/internal/testing"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepo(t *testing.T) {
	m := &testing2.RedisClientMock{}

	repo := redisRepo{
		client:  m,
		context: context.Background(),
	}

	sub := New(hash.Hash("from"), hash.Hash("to"), "sub")

	m.Queue("exists", int64(1), nil)
	assert.True(t, repo.Has(&sub))

	m.Queue("exists", int64(0), nil)
	assert.False(t, repo.Has(&sub))

	m.Queue("exists", int64(33), nil)
	assert.True(t, repo.Has(&sub))

	m.Queue("exists", int64(0), errors.New("key not exist"))
	assert.False(t, repo.Has(&sub))

	m.Queue("set", "foobar", nil)
	assert.NoError(t, repo.Store(&sub))

	m.Queue("set", "foobar", errors.New("error"))
	assert.Error(t, repo.Store(&sub))

	m.Queue("del", int64(1), nil)
	assert.NoError(t, repo.Remove(&sub))

	m.Queue("del", int64(0), errors.New("error"))
	assert.Error(t, repo.Remove(&sub))
}
