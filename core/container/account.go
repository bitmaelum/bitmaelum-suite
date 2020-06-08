package container

import (
    "github.com/jaytaph/mailv2/core/account"
    "github.com/jaytaph/mailv2/core/config"
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

    p, _ := homedir.Expand(config.Server.Account.Path)
    repo := account.NewFileRepository(p)
    accountRepository = &repo
    return accountRepository
}
