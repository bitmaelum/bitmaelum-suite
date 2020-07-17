package container

import (
	"github.com/bitmaelum/bitmaelum-suite/core/account"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
)

var accountService *account.Service

// GetAccountService retrieves an account service
func GetAccountService() *account.Service {
	if accountService != nil {
		return accountService
	}

	repo := account.NewFileRepository(config.Server.Paths.Accounts)

	accountService = account.NewService(repo)
	return accountService
}
