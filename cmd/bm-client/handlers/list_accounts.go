package handlers

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/pkg/vault"
	"github.com/olekukonko/tablewriter"
	"os"
)

// ListAccounts displays the current accounts available in the vault
func ListAccounts(vault *vault.Vault, displayKeys bool) {
	table := tablewriter.NewWriter(os.Stdout)

	headers := []string{"Default", "Address", "Name", "Server"}
	align := []int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT}

	if displayKeys {
		headers = append(headers, "Private Key", "Public Key")
		align = append(align, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT)
	}

	// alignment must be set at once
	table.SetColumnAlignment(align)
	table.SetHeader(headers)

	for _, acc := range vault.Accounts {
		def := ""
		if acc.Default {
			def = "*"
		}
		values := []string{
			def,
			acc.Address,
			acc.Name,
			acc.Routing,
		}
		if displayKeys {
			values = append(values, acc.PrivKey.S, acc.PubKey.S)
		}

		table.Append(values)
	}
	table.Render() // Send output
}
