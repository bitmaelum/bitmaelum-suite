package container

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/invite"
	"github.com/mitchellh/go-homedir"

	"github.com/go-redis/redis/v8"
)

// GetInviteRepo retrieves an invitation service
func GetInviteRepo() invite.Repository {

	//If redis.host is set on the config file it will use redis instead of bolt
	if config.Server.Redis.Host != "" {
		opts := redis.Options{
			Addr: config.Server.Redis.Host,
			DB:   config.Server.Redis.Db,
		}

		return invite.NewRedisRepository(&opts)
	}

	//If redis is not set then it will use BoltDB as default
	dbPath, _ := homedir.Expand(config.Server.Bolt.DatabasePath)
	return invite.NewBoltRepository(&dbPath)
}
