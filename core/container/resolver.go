package container

import (
	"github.com/bitmaelum/bitmaelum-suite/core/resolve"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
)

// We can have multiple resolvers to resolve a single address. We could resolve locally, remotely through resolver-services, or through DHT.
// We chain them all together with the ChainRepository

var (
	resolveService           *resolve.Service
	localResolverRepository  *resolve.Repository
	remoteResolverRepository *resolve.Repository
	dhtResolverRepository    *resolve.Repository
	chainResolverRepository  *resolve.ChainRepository
)

// GetResolveService retrieves a resolver service
func GetResolveService() *resolve.Service {
	if resolveService != nil {
		return resolveService
	}

	repo := getChainRepository()
	_ = repo.Add(*getLocalRepository())
	_ = repo.Add(*getRemoteRepository())
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

func getLocalRepository() *resolve.Repository {
	if localResolverRepository != nil {
		return localResolverRepository
	}

	repo := resolve.NewLocalRepository(GetAccountService())
	localResolverRepository = &repo
	return localResolverRepository
}

func getRemoteRepository() *resolve.Repository {
	if remoteResolverRepository != nil {
		return remoteResolverRepository
	}

	repo := resolve.NewRemoteRepository(config.Server.Resolver.Remote.URL)
	remoteResolverRepository = &repo
	return remoteResolverRepository
}

func getDhtRepository() *resolve.Repository {
	if dhtResolverRepository != nil {
		return dhtResolverRepository
	}

	repo := resolve.NewDHTRepository()
	dhtResolverRepository = &repo
	return dhtResolverRepository
}
