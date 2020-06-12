package container

import (
    "github.com/bitmaelum/bitmaelum-server/core/account"
    "github.com/bitmaelum/bitmaelum-server/core/config"
    "github.com/mitchellh/go-homedir"
)

var accountService *account.Service = nil
var accountRepository *account.Repository = nil

func GetAccountService() *account.Service{
    if accountService != nil {
        return accountService;
    }

    repo := GetAccountRepository()
    accountService = account.AccountService(*repo)
    return accountService
}

func GetAccountRepository() *account.Repository {
    if accountRepository != nil {
        return accountRepository;
    }

    p, _ := homedir.Expand(config.Server.Accounts.Path)
    repo := account.NewFileRepository(p)
    accountRepository = &repo
    return accountRepository
}
