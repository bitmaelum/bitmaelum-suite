package main

import (
	"os"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/olekukonko/tablewriter"
)

type options struct {
	Config   string `short:"c" long:"config" description:"Path to your configuration file"`
	Password string `short:"p" long:"password" description:"Vault password" default:""`
}

var opts options

func main() {
	internal.ParseOptions(&opts)
	config.LoadClientConfig(opts.Config)

	vault.VaultPassword = opts.Password
	v := vault.OpenVault()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Address", "New"})

	for _, acc := range v.Data.Accounts {
		addr, err := address.New(acc.Address)
		if err != nil {
			panic(err)
		}

		values := []string{
			addr.String(),
			addr.Hash().String(),
		}

		table.Append(values)
	}

	table.Render()
}
