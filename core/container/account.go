package container

import (
	"github.com/bitmaelum/bitmaelum-suite/core/account"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/mitchellh/go-homedir"
)

var accountService *account.Service

// GetAccountService retrieves an account service
func GetAccountService() *account.Service {
	if accountService != nil {
		return accountService
	}

	p, _ := homedir.Expand(config.Server.Paths.Accounts)
	repo := account.NewFileRepository(p)

	accountService = account.NewService(repo)
	return accountService
}
