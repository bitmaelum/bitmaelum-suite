package container

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/subscription"
	"github.com/go-redis/redis/v8"
)

// GetSubscriptionRepo returns the repository for storing and fetching subscriptions
func GetSubscriptionRepo() subscription.Repository {
	opts := redis.Options{
		Addr: config.Server.Redis.Host,
		DB:   config.Server.Redis.Db,
	}

	return subscription.NewRepository(&opts)
}
