package container

import (
    "github.com/go-redis/redis/v8"
    "github.com/jaytaph/mailv2/core/config"
    "github.com/jaytaph/mailv2/server/incoming"
    logger "github.com/sirupsen/logrus"
)

var incomingService *incoming.Service = nil
var incomingRepository *incoming.Repository = nil

func GetIncomingService() *incoming.Service{
    if incomingService != nil {
        logger.Trace("Returning cached incomingService")
        return incomingService;
    }

    logger.Trace("Creating new incomingService")
    repo := GetIncomingRepository()
    incomingService = incoming.NewIncomingService(*repo)
    return incomingService
}

func GetIncomingRepository() *incoming.Repository {
    if incomingRepository != nil {
        logger.Trace("Returning cached incomingRepository")
        return incomingRepository;
    }

    logger.Trace("Creating new incomingRepository")
    opts := redis.Options{
        Addr: config.Configuration.Redis.Host,
        DB: config.Configuration.Redis.Db,
    }

    repo := incoming.NewRedisRepository(&opts)
    incomingRepository = &repo
    return incomingRepository
}
