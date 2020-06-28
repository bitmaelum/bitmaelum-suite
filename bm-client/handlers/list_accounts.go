package handlers

import (
	"github.com/bitmaelum/bitmaelum-server/bm-client/account"
	"github.com/olekukonko/tablewriter"
	"os"
)

func ListAccounts(displayKeys bool) {
	table := tablewriter.NewWriter(os.Stdout)

	headers := []string{"Address", "Name", "Organisation", "Server"}
	if displayKeys {
		headers = append(headers, "Private Key", "Public Key")
	}
	table.SetHeader(headers)

	for _, acc := range account.Vault.Accounts {
		values := []string{
			acc.Address,
			acc.Name,
			acc.Organisation,
			acc.Server,
		}
		if displayKeys {
			values = append(values, acc.PrivKey, acc.PubKey)
		}

		table.Append(values)
	}
	table.Render() // Send output
}
