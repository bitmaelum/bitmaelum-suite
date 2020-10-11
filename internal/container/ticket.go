package container

import (
	"sync"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"

	"github.com/go-redis/redis/v8"
)

var (
	ticketOnce       sync.Once
	ticketRepository ticket.Repository
)

// GetTicketRepo returns the repository for storing and fetching tickets
func GetTicketRepo() ticket.Repository {

	ticketOnce.Do(func() {
		//If redis.host is set on the config file it will use redis instead of bolt
		if config.Server.Redis.Host != "" {
			opts := redis.Options{
				Addr: config.Server.Redis.Host,
				DB:   config.Server.Redis.Db,
			}

			ticketRepository = ticket.NewRedisRepository(&opts)
			return
		}

		//If redis is not set then it will use BoltDB as default
		ticketRepository = ticket.NewBoltRepository(config.Server.Bolt.DatabasePath)
	})

	return ticketRepository
}
