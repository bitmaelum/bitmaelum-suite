package handlers

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/vault"
	"github.com/olekukonko/tablewriter"
	"os"
)

// ListAccounts displays the current accounts available in the vault
func ListAccounts(vault *vault.Vault, displayKeys bool) {
	table := tablewriter.NewWriter(os.Stdout)

	headers := []string{"Address", "Name", "Organisation", "Server"}
	if displayKeys {
		headers = append(headers, "Private Key", "Public Key")
	}
	table.SetHeader(headers)

	for _, acc := range vault.Accounts {
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
