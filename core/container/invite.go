package container

import (
	"github.com/bitmaelum/bitmaelum-suite/core/invite"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/go-redis/redis/v8"
)

var inviteService *invite.Service

// GetInviteService retrieves an invitation service
func GetInviteService() *invite.Service {
	if inviteService != nil {
		return inviteService
	}

	opts := redis.Options{
		Addr: config.Server.Redis.Host,
		DB:   config.Server.Redis.Db,
	}

	repo := invite.NewRedisRepository(&opts)
	inviteService = invite.NewInviteService(repo)
	return inviteService
}
