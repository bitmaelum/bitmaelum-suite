// Copyright (c) 2021 BitMaelum Authors
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

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listOrganisationCmd = &cobra.Command{
	Use:   "list",
	Short: "List organisations",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenDefaultVault()

		if *oaddress == "" {
			displayOrganisations(v)
		} else {
			orgHash := hash.New(*oaddress)

			info, err := v.GetOrganisationInfo(orgHash)
			if err != nil {
				fmt.Println("error: organisation not found: ", *oaddress)
				os.Exit(1)
			}

			displayOrganisation(*info)
		}
	},
}

func displayOrganisations(v *vault.Vault) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Organisation", "Name", "Validations"})

	for _, o := range v.Store.Organisations {
		addr := "...@" + o.Addr + "!"

		table.Append([]string{
			addr,
			o.FullName,
			strings.Join(validations(o), "\n"),
		})

	}

	table.Render()
}

func validations(org vault.OrganisationInfo) []string {
	var ret []string

	for _, val := range org.Validations {
		var valStr string
		if ok, err := val.Validate(*org.ToOrg()); err == nil && ok {
			valStr = "[\U00002713] " + val.String()
		} else {
			valStr = "[\U00002717] " + val.String()
		}

		ret = append(ret, valStr)
	}

	return ret
}

func displayOrganisation(info vault.OrganisationInfo) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorders(tablewriter.Border{Left: false, Top: true, Right: false, Bottom: true})
	table.SetCenterSeparator("|")
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)

	table.AppendBulk([][]string{
		{"Organisation", "...@" + info.Addr + "!"},
		{"Organisation hash", info.ToOrg().Hash.String()},
	})

	table.AppendBulk([][]string{
		{"", ""},
		{"Full name", info.FullName},
		{"", ""},
		{"Public key", strings.Join(chunks(info.PubKey.String(), 78), "\n")},
		{"Proof of work", fmt.Sprintf("%d bits", info.Pow.Bits)},
	})

	if len(info.Validations) > 0 {
		table.AppendBulk([][]string{
			{"", ""},
			{"Validations", strings.Join(validations(info), "\n")},
		})
	}

	table.Render()
}

var (
	oaddress *string
)

func init() {
	organisationCmd.AddCommand(listOrganisationCmd)

	oaddress = listOrganisationCmd.Flags().String("organisation", "", "Organisation to display")
}
