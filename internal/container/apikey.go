package container

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/apikey"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/go-redis/redis/v8"
)

// GetAPIKeyRepo returns the repository for storing and fetching api keys
func GetAPIKeyRepo() apikey.Repository {
	opts := redis.Options{
		Addr: config.Server.Redis.Host,
		DB:   config.Server.Redis.Db,
	}

	return apikey.NewRepository(&opts)
}
