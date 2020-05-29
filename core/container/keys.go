package container

import (
    "github.com/jaytaph/mailv2/core/config"
    "github.com/jaytaph/mailv2/core/resolve"
)

var keysService *resolve.Service = nil

var localKeysRepository *resolve.Repository = nil

var remoteKeysRepository *resolve.Repository = nil

var dhtKeysRepository *resolve.Repository = nil

var chainKeysRepository *resolve.ChainRepository = nil


func GetKeyRetrievalService() *resolve.Service{
    if keysService != nil {
        return keysService;
    }

    repo := getChainRepository()
    repo.Add(*getLocalRepository())
    repo.Add(*getRemoteRepository())
    repo.Add(*getDhtRepository())

    keysService = resolve.KeyRetrievalService(repo)
    return keysService
}

func getChainRepository() *resolve.ChainRepository {
    if chainKeysRepository != nil {
        return chainKeysRepository;
    }

    chainKeysRepository = resolve.NewChainRepository()
    return chainKeysRepository
}

func getLocalRepository() *resolve.Repository {
    if localKeysRepository != nil {
        return localKeysRepository;
    }

    repo := resolve.NewLocalRepository(GetAccountService())
    localKeysRepository = &repo
    return localKeysRepository
}

func getRemoteRepository() *resolve.Repository {
    if remoteKeysRepository != nil {
        return remoteKeysRepository;
    }

    repo := resolve.NewRemoteRepository(config.Client.Resolve.Remote.Url)
    remoteKeysRepository = &repo
    return remoteKeysRepository
}

func getDhtRepository() *resolve.Repository {
    if dhtKeysRepository != nil {
        return dhtKeysRepository;
    }

    repo := resolve.NewDHTRepository()
    dhtKeysRepository = &repo
    return dhtKeysRepository
}
