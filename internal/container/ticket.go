package container

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"

	"github.com/go-redis/redis/v8"
)

// GetTicketRepo returns the repository for storing and fetching tickets
func GetTicketRepo() ticket.Repository {

	//If redis.host is set on the config file it will use redis instead of bolt
	if config.Server.Redis.Host != "" {
		opts := redis.Options{
			Addr: config.Server.Redis.Host,
			DB:   config.Server.Redis.Db,
		}

		return ticket.NewRedisRepository(&opts)
	}

	//If redis is not set then it will use BoltDB as default
	return ticket.NewBoltRepository(&config.Server.Bolt.DatabasePath)
}
