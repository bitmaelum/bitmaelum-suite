package container

import (
	"github.com/bitmaelum/bitmaelum-server/core/resolve"
	"github.com/bitmaelum/bitmaelum-server/internal/config"
)

var (
	resolveService       *resolve.Service
	localKeysRepository  *resolve.Repository
	remoteKeysRepository *resolve.Repository
	dhtKeysRepository    *resolve.Repository
	chainKeysRepository  *resolve.ChainRepository
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
	if chainKeysRepository != nil {
		return chainKeysRepository
	}

	chainKeysRepository = resolve.NewChainRepository()
	return chainKeysRepository
}

func getLocalRepository() *resolve.Repository {
	if localKeysRepository != nil {
		return localKeysRepository
	}

	repo := resolve.NewLocalRepository(GetAccountService())
	localKeysRepository = &repo
	return localKeysRepository
}

func getRemoteRepository() *resolve.Repository {
	if remoteKeysRepository != nil {
		return remoteKeysRepository
	}

	repo := resolve.NewRemoteRepository(config.Client.Resolver.Remote.URL)
	remoteKeysRepository = &repo
	return remoteKeysRepository
}

func getDhtRepository() *resolve.Repository {
	if dhtKeysRepository != nil {
		return dhtKeysRepository
	}

	repo := resolve.NewDHTRepository()
	dhtKeysRepository = &repo
	return dhtKeysRepository
}
