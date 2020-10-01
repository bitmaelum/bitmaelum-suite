package handlers

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/olekukonko/tablewriter"
	"os"
)

// ListAccounts displays the current accounts available in the vault
func ListAccounts(vault *vault.Vault, displayKeys bool) {
	table := tablewriter.NewWriter(os.Stdout)

	headers := []string{"Address", "Name", "Routing ID", "Route Host"}
	align := []int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT}

	if displayKeys {
		headers = append(headers, "Private Key", "Public Key")
		align = append(align, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT)
	}

	// alignment must be set at once
	table.SetColumnAlignment(align)
	table.SetHeader(headers)

	for _, acc := range vault.Data.Accounts {
		if acc.Default {
			acc.Address += " (*)"
		}

		var route string
		resolver := container.GetResolveService()
		r, err := resolver.ResolveRouting(acc.RoutingID)
		if err != nil {
			route = "route not found"
		} else {
			route = r.Routing
		}

		values := []string{
			acc.Address,
			acc.Name,
			acc.RoutingID[:10],
			route,
		}
		if displayKeys {
			values = append(values, acc.PrivKey.S, acc.PubKey.S)
		}

		table.Append(values)
	}
	table.Render() // Send output
}
