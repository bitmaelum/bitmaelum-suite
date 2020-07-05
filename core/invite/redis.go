package invite

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/bitmaelum/bitmaelum-server/pkg/address"
	"github.com/go-redis/redis/v8"
	"time"
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

// CreateInvite generate a new invitation and stores this in redis
func (r *redisRepo) CreateInvite(addr address.HashAddress, expiry time.Duration) (string, error) {
	buff := make([]byte, 32)
	_, err := rand.Read(buff)
	if err != nil {
		return "", err
	}
	token := base64.StdEncoding.EncodeToString(buff)

	err = r.client.Set(r.client.Context(), "invite."+addr.String(), token, expiry).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetInvite retrieves an invite from redis
func (r *redisRepo) GetInvite(addr address.HashAddress) (string, error) {
	return r.client.Get(r.client.Context(), "invite."+addr.String()).Result()
}

// RemoveInvite deletes an invite from redis
func (r *redisRepo) RemoveInvite(addr address.HashAddress) error {
	return r.client.Del(r.client.Context(), "invite."+addr.String()).Err()
}
