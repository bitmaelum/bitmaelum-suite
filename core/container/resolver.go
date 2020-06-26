package container

import (
	"github.com/bitmaelum/bitmaelum-server/core/config"
	"github.com/bitmaelum/bitmaelum-server/core/resolve"
)

var resolveService *resolve.Service = nil

var localKeysRepository *resolve.Repository = nil

var remoteKeysRepository *resolve.Repository = nil

var dhtKeysRepository *resolve.Repository = nil

var chainKeysRepository *resolve.ChainRepository = nil

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

	repo := resolve.NewRemoteRepository(config.Client.Resolver.Remote.Url)
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
