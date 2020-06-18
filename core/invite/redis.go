package invite

import (
    "crypto/rand"
    "encoding/base64"
    "github.com/bitmaelum/bitmaelum-server/core"
    "github.com/go-redis/redis/v8"
    "time"
)

type redisRepo struct {
    client *redis.Client
}

func NewRedisRepository(opts *redis.Options) Repository {
    return &redisRepo{
        client: redis.NewClient(opts),
    }
}

func (r *redisRepo) CreateInvite(addr core.HashAddress, expiry time.Duration) (string, error) {
    buff := make([]byte, 32)
    _, err := rand.Read(buff)
    if err != nil {
        return "", err
    }
    token := base64.StdEncoding.EncodeToString(buff)

    err = r.client.Set(r.client.Context(), "invite." + addr.String(), token, expiry).Err()
    if err != nil {
        return "", err
    }

    return token, nil
}

func (r *redisRepo) GetInvite(addr core.HashAddress) (string, error) {
    return r.client.Get(r.client.Context(), "invite." + addr.String()).Result()
}

func (r *redisRepo) RemoveInvite(addr core.HashAddress) error {
    return r.client.Del(r.client.Context(), "invite." + addr.String()).Err()
}
