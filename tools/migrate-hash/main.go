package main

import (
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	tools "github.com/bitmaelum/bitmaelum-suite/tools/internal"
	"github.com/olekukonko/tablewriter"
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

	tools.VaultPassword = opts.Password
	v := tools.OpenVault()


	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Address", "Old", "New"})

	for _, acc := range v.Accounts {
		addr, err := address.New(acc.Address)
		if err != nil {
			panic(err)
		}

		values := []string{
			addr.String(),
			addr.OldHash().String(),
			addr.Hash().String(),
		}

		table.Append(values)
	}

	table.Render()
}
