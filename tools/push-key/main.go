package main

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-server/core"
	"github.com/bitmaelum/bitmaelum-server/core/container"
	"github.com/bitmaelum/bitmaelum-server/core/password"
	"os"
)

type options struct {
	Config   string `short:"c" long:"config" description:"Path to your configuration file"`
	Password string `short:"p" long:"password" description:"Vault password" default:""`
}

var opts options

func main() {
	core.ParseOptions(&opts)
	core.LoadClientConfig(opts.Config)

	if opts.Password == "" {
		opts.Password = string(password.AskPassword())
	}

	// Unlock vault
	err := account.UnlockVault(config.Client.Accounts.Path, []byte(opts.Password))
	if err != nil {
		fmt.Printf("Error while opening vault: %s", err)
		fmt.Println("")
		os.Exit(1)
	}

	rs := container.GetResolveService()
	for i, acc := range account.Vault.Accounts {
		account.Vault.Accounts[i].Server = "bitmaelum.ngrok.io:443"
		account.Vault.Save()

		err = rs.UploadInfo(acc, acc.Server)
		if err != nil {
			fmt.Printf("Error for account %s: %s\n", acc.Address, err.Error())
		}
	}
}
