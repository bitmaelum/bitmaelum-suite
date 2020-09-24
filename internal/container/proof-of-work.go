package container

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/storage"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/go-homedir"
)

var powService storage.Storable

// GetProofOfWorkService returns a service that can store a proof of work
func GetProofOfWorkService() storage.Storable {
	if powService != nil {
		return powService
	}

	//If redis.host is set on the config file it will use redis instead of bolt
	if config.Server.Redis.Host != "" {
		opts := redis.Options{
			Addr: config.Server.Redis.Host,
			DB:   config.Server.Redis.Db,
		}

		powService = storage.NewRedis(&opts)
	} else {
		dbPath, _ := homedir.Expand(config.Server.Bolt.DatabasePath)
		powService = storage.NewBolt(&dbPath)
	}

	return powService
}
