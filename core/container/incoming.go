package container

import (
    "github.com/go-redis/redis/v8"
    "github.com/jaytaph/mailv2/core/config"
    "github.com/jaytaph/mailv2/server/incoming"
)

var incomingService *incoming.Service = nil
var incomingRepository *incoming.Repository = nil

func GetIncomingService() *incoming.Service{
    if incomingService != nil {
        return incomingService;
    }

    repo := GetIncomingRepository()
    incomingService = incoming.NewIncomingService(*repo)
    return incomingService
}

func GetIncomingRepository() *incoming.Repository {
    if incomingRepository != nil {
        return incomingRepository;
    }

    opts := redis.Options{
        Addr: config.Configuration.Redis.Host,
        DB: config.Configuration.Redis.Db,
    }

    repo := incoming.NewRedisRepository(&opts)
    incomingRepository = &repo
    return incomingRepository
}
