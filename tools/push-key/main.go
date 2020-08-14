package main

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/vault"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/password"
	"os"
)

type options struct {
	Config   string `short:"c" long:"config" description:"Path to your configuration file"`
	Password string `short:"p" long:"password" description:"Vault password" default:""`
}

var opts options

func main() {
	internal.ParseOptions(&opts)
	config.LoadClientConfig(opts.Config)
	config.LoadServerConfig(opts.Config)

	fromVault := false
	if opts.Password == "" {
		opts.Password, fromVault = password.AskPassword()
	}

	// Unlock vault
	accountVault, err := vault.New(config.Client.Accounts.Path, []byte(opts.Password))
	if err != nil {
		fmt.Printf("Error while opening vault: %s", err)
		fmt.Println("")
		os.Exit(1)
	}

	// If the password was correct and not already read from the vault, store it in the vault
	if !fromVault {
		_ = password.StorePassword(opts.Password)
	}

	rs := container.GetResolveService()
	for _, acc := range accountVault.Accounts {
		// accountVault.Accounts[i].Server = "bitmaelum.ngrok.io:443"
		// accountVault.Save()

		err = rs.UploadInfo(acc, acc.Server)
		if err != nil {
			fmt.Printf("Error for account %s: %s\n", acc.Address, err.Error())
		}
	}
}
