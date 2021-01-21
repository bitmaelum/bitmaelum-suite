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

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/organisation"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/spf13/cobra"
)

var organisationValidateRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an organisation validation",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenDefaultVault()

		orgHash := hash.New(strings.TrimRight(*ovOrganisation, "!"))

		info, err := v.GetOrganisationInfo(orgHash)
		if err != nil {
			fmt.Println("error: organisation not found: ", *oaddress)
			os.Exit(1)
		}

		val, err := organisation.NewValidationTypeFromString(fmt.Sprintf("%s %s", *ovrType, *ovrValue))
		if err != nil {
			fmt.Println("error: incorrect validation type/value: ", err)
			os.Exit(1)
		}

		newValidations := []organisation.ValidationType{}
		for _, srcVal := range info.Validations {
			if srcVal.String() != val.String() {
				newValidations = append(newValidations, srcVal)
			}
		}
		info.Validations = newValidations

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

		fmt.Println("Successfully removed validation")
	},
}

var (
	ovrType  *string
	ovrValue *string
)

func init() {
	organisationValidateCmd.AddCommand(organisationValidateRemoveCmd)

	ovrType = organisationValidateRemoveCmd.PersistentFlags().String("type", "", "Type")
	ovrValue = organisationValidateRemoveCmd.PersistentFlags().String("value", "", "Value")

	_ = organisationValidateRemoveCmd.MarkPersistentFlagRequired("type")
	_ = organisationValidateRemoveCmd.MarkPersistentFlagRequired("value")
}
