package container

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/invite"
	"github.com/go-redis/redis/v8"
)

// GetInviteRepo retrieves an invitation service
func GetInviteRepo() invite.Repository {
	opts := redis.Options{
		Addr: config.Server.Redis.Host,
		DB:   config.Server.Redis.Db,
	}

	return invite.NewRedisRepository(&opts)
}
