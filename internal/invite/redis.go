package invite

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
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

// Create generate a new invitation and stores this in redis
func (r *redisRepo) Create(addr address.HashAddress, expiry time.Duration) (string, error) {
	buff := make([]byte, 32)
	_, err := rand.Read(buff)
	if err != nil {
		return "", err
	}
	token := base64.StdEncoding.EncodeToString(buff)

	err = r.client.Set(r.client.Context(), createInviteKey(addr), token, expiry).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

// Get retrieves an invite from redis
func (r *redisRepo) Get(addr address.HashAddress) (string, error) {
	return r.client.Get(r.client.Context(), createInviteKey(addr)).Result()
}

// Remove deletes an invite from redis
func (r *redisRepo) Remove(addr address.HashAddress) error {
	return r.client.Del(r.client.Context(), createInviteKey(addr)).Err()
}
