package container

import (
	"sync"

	"github.com/bitmaelum/bitmaelum-suite/internal/apikey"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/go-redis/redis/v8"
)

var (
	apikeyOnce       sync.Once
	apikeyRepository apikey.Repository
)

// GetAPIKeyRepo returns the repository for storing and fetching api keys
func GetAPIKeyRepo() apikey.Repository {

	apikeyOnce.Do(func() {
		//If redis.host is set on the config file it will use redis instead of bolt
		if config.Server.Redis.Host != "" {
			opts := redis.Options{
				Addr: config.Server.Redis.Host,
				DB:   config.Server.Redis.Db,
			}

			apikeyRepository = apikey.NewRedisRepository(&opts)
			return
		}

		//If redis is not set then it will use BoltDB as default
		apikeyRepository = apikey.NewBoltRepository(&config.Server.Bolt.DatabasePath)
	})

	return apikeyRepository
}
