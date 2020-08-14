package container

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"
	"github.com/go-redis/redis/v8"
)

// GetTicketRepo returns the repository for storing and fetching tickets
func GetTicketRepo() ticket.Repository {
	opts := redis.Options{
		Addr: config.Server.Redis.Host,
		DB:   config.Server.Redis.Db,
	}

	return ticket.NewRepository(&opts)
}
