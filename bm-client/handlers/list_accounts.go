package handlers

import (
	"github.com/bitmaelum/bitmaelum-server/bm-client/account"
	"github.com/olekukonko/tablewriter"
	"os"
)

func ListAccounts() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Address", "Name", "Organisation", "Server"})

	for _, acc := range account.Vault.Accounts {
		table.Append([]string{
			acc.Address,
			acc.Name,
			acc.Organisation,
			acc.Server,
		})
	}
	table.Render() // Send output
}
