package internal

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/pkg/vault"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/password"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"
	"os"
)

// VaultPassword is the given password through the commandline for opening the vault
var VaultPassword string

// OpenVault returns an opened vault, or opens the vault, asking a password if needed
func OpenVault() *vault.Vault {
	fromVault := false
	if VaultPassword == "" {
		VaultPassword, fromVault = password.AskPassword()
	}

	// Unlock vault
	v, err := vault.New(config.Client.Accounts.Path, []byte(VaultPassword))
	if err != nil {
		fmt.Printf("Error while opening vault: %s", err)
		fmt.Println("")
		os.Exit(1)
	}

	// If the password was correct and not already read from the vault, store it in the vault
	if !fromVault {
		_ = password.StorePassword(VaultPassword)
	}

	return v
}

// GetAccountOrDefault find the address from the vault. If address is empty, it will fetch the default address, or the
// first address in the vault if no default address is present.
func GetAccountOrDefault(vault *vault.Vault, a string) *internal.AccountInfo {
	if a == "" {
		return vault.GetDefaultAccount()
	}

	addr, err := address.New(a)
	if err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}

	info, err := vault.GetAccountInfo(*addr)
	if err != nil {
		logrus.Fatal("Address not found in vault")
		os.Exit(1)
	}

	return info
}
