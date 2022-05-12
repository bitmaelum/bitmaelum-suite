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

var organisationValidateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List organisation validations",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenDefaultVault()

		orgHash := hash.New(strings.TrimRight(*ovOrganisation, "!"))

		info, err := v.GetOrganisationInfo(orgHash)
		if err != nil {
			fmt.Println("error: organisation not found: ", *oaddress)
			os.Exit(1)
		}

		displayValidations(*info)
	},
}

func displayValidations(org vault.OrganisationInfo) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Type", "Value", "Validated"})

	for _, v := range org.Validations {
		var valStr string

		if ok, err := v.Validate(*org.ToOrg()); err == nil && ok {
			valStr = "yes"
		} else {
			valStr = "no"
		}

		table.Append([]string{
			v.Type,
			v.Value,
			valStr,
		})

	}

	table.Render()
}

func init() {
	organisationValidateCmd.AddCommand(organisationValidateListCmd)
}
