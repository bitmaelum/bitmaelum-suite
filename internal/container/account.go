package container

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/account"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
)

// GetAccountRepo retrieves an account repository
func GetAccountRepo() account.Repository {
	return account.NewFileRepository(config.Server.Paths.Accounts)
}
