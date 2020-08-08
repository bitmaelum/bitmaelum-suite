package container

import (
	account2 "github.com/bitmaelum/bitmaelum-suite/internal/account"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
)

// GetAccountRepo retrieves an account repository
func GetAccountRepo() account2.Repository {
	return account2.NewFileRepository(config.Server.Paths.Accounts)
}
