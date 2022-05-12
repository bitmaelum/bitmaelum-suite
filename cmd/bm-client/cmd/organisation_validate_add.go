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

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/organisation"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/spf13/cobra"
)

var organisationValidateAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add organisation validation",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenDefaultVault()

		orgHash := hash.New(strings.TrimRight(*ovOrganisation, "!"))

		info, err := v.GetOrganisationInfo(orgHash)
		if err != nil {
			fmt.Println("error: organisation not found: ", *oaddress)
			os.Exit(1)
		}

		val, err := organisation.NewValidationTypeFromString(fmt.Sprintf("%s %s", *ovaType, *ovaValue))
		if err != nil {
			fmt.Println("error: incorrect validation type/value: ", err)
			os.Exit(1)
		}

		for i := range info.Validations {
			if info.Validations[i].Type == val.Type && info.Validations[i].Value == val.Value {
				fmt.Println("error: type/value already present")
				os.Exit(1)
			}
		}

		info.Validations = append(info.Validations, *val)

		err = v.Persist()
		if err != nil {
			fmt.Println("error: cannot save data back into the vault: ", err)
			os.Exit(1)
		}

		rs := container.Instance.GetResolveService()
		err = rs.UploadOrganisationInfo(*info)
		if err != nil {
			fmt.Println("error: cannot upload data to the resolver: ", err)
			os.Exit(1)
		}

		fmt.Println("Successfully added validation")
	},
}

var (
	ovaType  *string
	ovaValue *string
)

func init() {
	organisationValidateCmd.AddCommand(organisationValidateAddCmd)

	ovaType = organisationValidateAddCmd.PersistentFlags().String("type", "", "Type")
	ovaValue = organisationValidateAddCmd.PersistentFlags().String("value", "", "Value")

	_ = organisationValidateAddCmd.MarkPersistentFlagRequired("type")
	_ = organisationValidateAddCmd.MarkPersistentFlagRequired("value")
}
