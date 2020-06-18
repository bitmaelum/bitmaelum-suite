package core

import (
	"github.com/bitmaelum/bitmaelum-server/core/account/client"
	"github.com/olekukonko/tablewriter"
	"os"
)

func ListAccounts() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Address", "Name", "Organisation"})

	accounts, err := client.GetAccounts()
	if err != nil {
		panic(err)
	}

	for _, account := range accounts {
		table.Append([]string{
			account.Address,
			account.Name,
			account.Organisation,
		})
	}
	table.Render() // Send output
}