package main

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-server/cmd/bm-client/vault"
	"github.com/bitmaelum/bitmaelum-server/core"
	"github.com/bitmaelum/bitmaelum-server/core/container"
	"github.com/bitmaelum/bitmaelum-server/core/password"
	"github.com/bitmaelum/bitmaelum-server/internal/config"
	"os"
)

type options struct {
	Config   string `short:"c" long:"config" description:"Path to your configuration file"`
	Password string `short:"p" long:"password" description:"Vault password" default:""`
}

var opts options

func main() {
	core.ParseOptions(&opts)
	config.LoadClientConfig(opts.Config)

	if opts.Password == "" {
		opts.Password = string(password.AskPassword())
	}

	// Unlock vault
	accountVault, err := vault.New(config.Client.Accounts.Path, []byte(opts.Password))
	if err != nil {
		fmt.Printf("Error while opening vault: %s", err)
		fmt.Println("")
		os.Exit(1)
	}

	rs := container.GetResolveService()
	for i, acc := range accountVault.Accounts {
		accountVault.Accounts[i].Server = "bitmaelum.ngrok.io:443"
		accountVault.Save()

		err = rs.UploadInfo(acc, acc.Server)
		if err != nil {
			fmt.Printf("Error for account %s: %s\n", acc.Address, err.Error())
		}
	}
}
