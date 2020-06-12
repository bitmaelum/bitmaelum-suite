package container

import (
    "github.com/go-redis/redis/v8"
    "github.com/bitmaelum/bitmaelum-server/core/config"
    "github.com/bitmaelum/bitmaelum-server/core/invite"
)

var inviteService *invite.Service = nil
var inviteRepository *invite.Repository = nil

func GetInviteService() *invite.Service{
    if inviteService != nil {
        return inviteService;
    }

    repo := GetInviteRepository()
    inviteService = invite.NewInviteService(*repo)
    return inviteService
}

func GetInviteRepository() *invite.Repository {
    if inviteRepository != nil {
        return inviteRepository;
    }

    opts := redis.Options{
        Addr: config.Server.Redis.Host,
        DB: config.Server.Redis.Db,
    }

    repo := invite.NewRedisRepository(&opts)
    inviteRepository = &repo
    return inviteRepository
}
