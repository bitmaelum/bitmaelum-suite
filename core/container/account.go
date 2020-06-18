package container

import (
	"github.com/bitmaelum/bitmaelum-server/core/account/server"
	"github.com/bitmaelum/bitmaelum-server/core/config"
	"github.com/mitchellh/go-homedir"
)

var accountService *server.Service = nil
var accountRepository *server.Repository = nil

func GetAccountService() *server.Service {
    if accountService != nil {
        return accountService;
    }

    repo := GetAccountRepository()
    accountService = server.AccountService(*repo)
    return accountService
}

func GetAccountRepository() *server.Repository {
    if accountRepository != nil {
        return accountRepository;
    }

    p, _ := homedir.Expand(config.Server.Accounts.Path)
    repo := server.NewFileRepository(p)
    accountRepository = &repo
    return accountRepository
}
