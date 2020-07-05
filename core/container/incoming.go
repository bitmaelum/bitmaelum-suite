package container

import (
	"github.com/bitmaelum/bitmaelum-server/cmd/bm-server/incoming"
	"github.com/bitmaelum/bitmaelum-server/internal/config"
	"github.com/go-redis/redis/v8"
)

var incomingService *incoming.Service
var incomingRepository *incoming.Repository

// GetIncomingService retrieves an incoming service
func GetIncomingService() *incoming.Service {
	if incomingService != nil {
		return incomingService
	}

	repo := getIncomingRepository()
	incomingService = incoming.NewIncomingService(*repo)
	return incomingService
}

func getIncomingRepository() *incoming.Repository {
	if incomingRepository != nil {
		return incomingRepository
	}

	opts := redis.Options{
		Addr: config.Server.Redis.Host,
		DB:   config.Server.Redis.Db,
	}

	repo := incoming.NewRedisRepository(&opts)
	incomingRepository = &repo
	return incomingRepository
}
