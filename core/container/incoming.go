package container

import (
    "github.com/bitmaelum/bitmaelum-server/bm-server/incoming"
	"github.com/bitmaelum/bitmaelum-server/core/config"
	"github.com/go-redis/redis/v8"
)

var incomingService *incoming.Service = nil
var incomingRepository *incoming.Repository = nil

func GetIncomingService() *incoming.Service{
    if incomingService != nil {
		return incomingService
    }

    repo := GetIncomingRepository()
    incomingService = incoming.NewIncomingService(*repo)
    return incomingService
}

func GetIncomingRepository() *incoming.Repository {
    if incomingRepository != nil {
		return incomingRepository
    }

    opts := redis.Options{
        Addr: config.Server.Redis.Host,
        DB: config.Server.Redis.Db,
    }

    repo := incoming.NewRedisRepository(&opts)
    incomingRepository = &repo
    return incomingRepository
}
