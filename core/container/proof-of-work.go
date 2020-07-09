package container

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/storage"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/go-redis/redis/v8"
)

var powService storage.Storable

// GetProofOfWorkService retrieves an incoming service
func GetProofOfWorkService() storage.Storable {
	if powService != nil {
		return powService
	}

	opts := redis.Options{
		Addr: config.Server.Redis.Host,
		DB:   config.Server.Redis.Db,
	}

	powService = storage.NewRedis(&opts)
	return powService
}
