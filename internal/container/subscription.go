package container

import (
	"sync"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/subscription"
	"github.com/go-redis/redis/v8"
)

var (
	subscriptionOnce       sync.Once
	subscriptionRepository subscription.Repository
)

// GetSubscriptionRepo returns the repository for storing and fetching subscriptions
func GetSubscriptionRepo() subscription.Repository {

	subscriptionOnce.Do(func() {
		//If redis.host is set on the config file it will use redis instead of bolt
		if config.Server.Redis.Host != "" {
			opts := redis.Options{
				Addr: config.Server.Redis.Host,
				DB:   config.Server.Redis.Db,
			}

			subscriptionRepository = subscription.NewRedisRepository(&opts)
			return
		}

		//If redis is not set then it will use BoltDB as default
		subscriptionRepository = subscription.NewBoltRepository(&config.Server.Bolt.DatabasePath)
	})

	return subscriptionRepository
}
