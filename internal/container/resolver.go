package container

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolve"
)

// We can have multiple resolvers to resolve a single address. We could resolve locally, remotely through resolver-services, or through DHT.
// We chain them all together with the ChainRepository

var (
	resolveService           *resolve.Service
	sqliteResolverRepository *resolve.Repository
	dhtResolverRepository    *resolve.Repository
	chainResolverRepository  *resolve.ChainRepository
)

// GetResolveService retrieves a resolver service
func GetResolveService() *resolve.Service {
	if resolveService != nil {
		return resolveService
	}

	repo := getChainRepository()
	if config.Server.Resolver.Sqlite.Enabled {
		r, err := getSQLiteRepository(config.Server.Resolver.Sqlite.Dsn)
		if err == nil {
			_ = repo.Add(r)
		}
	}
	if config.Client.Resolver.Remote.Enabled {
		_ = repo.Add(*getRemoteRepository(config.Client.Resolver.Remote.URL))
	}
	if config.Server.Resolver.Remote.Enabled {
		_ = repo.Add(*getRemoteRepository(config.Server.Resolver.Remote.URL))
	}
	_ = repo.Add(*getDhtRepository())

	return resolve.KeyRetrievalService(repo)
}

func getChainRepository() *resolve.ChainRepository {
	if chainResolverRepository != nil {
		return chainResolverRepository
	}

	chainResolverRepository = resolve.NewChainRepository()
	return chainResolverRepository
}

func getRemoteRepository(url string) *resolve.Repository {
	repo := resolve.NewRemoteRepository(url)
	return &repo
}

func getSQLiteRepository(dsn string) (resolve.Repository, error) {
	repo, err := resolve.NewSqliteRepository(dsn)
	return repo, err
}

func getDhtRepository() *resolve.Repository {
	if dhtResolverRepository != nil {
		return dhtResolverRepository
	}

	repo := resolve.NewDHTRepository()
	dhtResolverRepository = &repo
	return dhtResolverRepository
}
