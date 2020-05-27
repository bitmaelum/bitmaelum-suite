package container

import (
    "github.com/jaytaph/mailv2/core/keys"
)

var keysService *keys.Service = nil

var localKeysRepository *keys.Repository = nil

var remoteKeysRepository *keys.Repository = nil

var dhtKeysRepository *keys.Repository = nil

var chainKeysRepository *keys.ChainRepository = nil


func GetKeyRetrievalService() *keys.Service{
    if keysService != nil {
        return keysService;
    }

    repo := getChainRepository()
    repo.Add(*getLocalRepository())
    repo.Add(*getRemoteRepository())
    repo.Add(*getDhtRepository())

    keysService = keys.KeyRetrievalService(repo)
    return keysService
}

func getChainRepository() *keys.ChainRepository {
    if chainKeysRepository != nil {
        return chainKeysRepository;
    }

    chainKeysRepository = keys.NewChainRepository()
    return chainKeysRepository
}

func getLocalRepository() *keys.Repository {
    if localKeysRepository != nil {
        return localKeysRepository;
    }

    repo := keys.NewLocalRepository(GetAccountService())
    localKeysRepository = &repo
    return localKeysRepository
}

func getRemoteRepository() *keys.Repository {
    if remoteKeysRepository != nil {
        return remoteKeysRepository;
    }

    repo := keys.NewRemoteRepository()
    remoteKeysRepository = &repo
    return remoteKeysRepository
}

func getDhtRepository() *keys.Repository {
    if dhtKeysRepository != nil {
        return dhtKeysRepository;
    }

    repo := keys.NewDHTRepository()
    dhtKeysRepository = &repo
    return dhtKeysRepository
}
