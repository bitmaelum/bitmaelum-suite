package container

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
)

// We can have multiple resolvers to resolve a single address. We could resolve locally, remotely through resolver-services, or through DHT.
// We chain them all together with the ChainRepository

var (
	resolveService          *resolver.Service
	chainResolverRepository *resolver.ChainRepository
)

// SetResolveService allows you to easily set your own resolve service. Used for unit testing
func SetResolveService(s *resolver.Service) {
	resolveService = s
}

// GetResolveService retrieves a resolver service
func GetResolveService() *resolver.Service {
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
		_ = repo.Add(*getRemoteRepository(config.Client.Resolver.Remote.URL, config.Client.Server.DebugHTTP))
	}
	if config.Server.Resolver.Remote.Enabled {
		_ = repo.Add(*getRemoteRepository(config.Server.Resolver.Remote.URL, false))
	}

	resolveService = resolver.KeyRetrievalService(repo)
	return resolveService
}

func getChainRepository() *resolver.ChainRepository {
	if chainResolverRepository != nil {
		return chainResolverRepository
	}

	chainResolverRepository = resolver.NewChainRepository()
	return chainResolverRepository
}

func getRemoteRepository(url string, debug bool) *resolver.Repository {
	repo := resolver.NewRemoteRepository(url, debug)
	return &repo
}

func getSQLiteRepository(dsn string) (resolver.Repository, error) {
	repo, err := resolver.NewSqliteRepository(dsn)
	return repo, err
}
