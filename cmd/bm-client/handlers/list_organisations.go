package handlers

import (
	"os"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/olekukonko/tablewriter"
)

// ListOrganisations displays the current accounts available in the vault
func ListOrganisations(vault *vault.Vault, displayKeys bool) {
	table := tablewriter.NewWriter(os.Stdout)

	headers := []string{"Organisation", "Name", "ID", "Validation"}
	align := []int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT}

	if displayKeys {
		headers = append(headers, "Private Key", "Public Key")
		align = append(align, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT)
	}

	// alignment must be set at once
	table.SetColumnAlignment(align)
	table.SetHeader(headers)

	for _, org := range vault.Data.Organisations {
		o, err := internal.InfoToOrg(org)
		if err != nil {
			continue
		}

		if len(org.Validations) == 0 {
			values := []string{
				"...@" + org.Addr + "!",
				org.FullName,
				"-",
				o.Hash.String(),
			}
			if displayKeys {
				values = append(values, org.PrivKey.S, org.PubKey.S)
			}

			table.Append(values)
		}

		for i, val := range org.Validations {
			var valstr string
			if ok, err := val.Validate(*o); err == nil && ok {
				valstr = "\U00002713 " + val.String()
			} else {
				valstr = "\U00002717 " + val.String()
			}

			var values []string
			if i == 0 {
				// First entry
				values = []string{
					"...@" + org.Addr + "!",
					org.FullName,
					valstr,
					o.Hash.String(),
				}
			} else {
				// Additional validation rows
				values = []string{
					"",
					"",
					valstr,
					"",
				}
			}

			if displayKeys {
				values = append(values, org.PrivKey.S, org.PubKey.S)
			}

			table.Append(values)
		}
	}
	table.Render() // Send output
}
