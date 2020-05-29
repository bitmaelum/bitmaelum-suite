package container

import (
    "github.com/jaytaph/mailv2/core/account"
    "github.com/jaytaph/mailv2/core/config"
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

    repo := account.NewFileRepository(config.Server.Account.Path)
    accountRepository = &repo
    return accountRepository
}
