package container

import (
	"github.com/bitmaelum/bitmaelum-server/core/account"
	"github.com/bitmaelum/bitmaelum-server/core/config"
	"github.com/mitchellh/go-homedir"
)

var accountService *account.Service
var accountRepository *account.Repository

// GetAccountService retrieves an account service
func GetAccountService() *account.Service {
	if accountService != nil {
		return accountService
	}

	repo := getAccountRepository()
	accountService = account.NewService(*repo)
	return accountService
}

func getAccountRepository() *account.Repository {
	if accountRepository != nil {
		return accountRepository
	}

	p, _ := homedir.Expand(config.Server.Accounts.Path)
	repo := account.NewFileRepository(p)
	accountRepository = &repo
	return accountRepository
}
