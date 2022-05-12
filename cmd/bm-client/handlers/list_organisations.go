// Copyright (c) 2022 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package handlers

import (
	"os"

	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/olekukonko/tablewriter"
)

// ListOrganisations displays the current accounts available in the vault
func ListOrganisations(v *vault.Vault, displayKeys bool) {
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

	for _, org := range v.Store.Organisations {

		if len(org.Validations) == 0 {
			values := []string{
				"...@" + org.Addr + "!",
				org.FullName,
				"-",
				org.ToOrg().Hash.String(),
			}
			if displayKeys {
				values = append(values, org.GetActiveKey().PrivKey.S, org.GetActiveKey().PubKey.S)
			}

			table.Append(values)
		}

		for i, val := range org.Validations {
			var valStr string
			if ok, err := val.Validate(*org.ToOrg()); err == nil && ok {
				valStr = "\U00002713 " + val.String()
			} else {
				valStr = "\U00002717 " + val.String()
			}

			var values []string
			if i == 0 {
				// First entry
				values = []string{
					"...@" + org.Addr + "!",
					org.FullName,
					valStr,
					org.ToOrg().Hash.String(),
				}
			} else {
				// Additional validation rows
				values = []string{
					"",
					"",
					valStr,
					"",
				}
			}

			if displayKeys {
				values = append(values, org.GetActiveKey().PrivKey.S, org.GetActiveKey().PubKey.S)
			}

			table.Append(values)
		}
	}
	table.Render() // Send output
}
