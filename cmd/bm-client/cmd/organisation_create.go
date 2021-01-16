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
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/handlers"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var organisationCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new organisation",
	Long: `Creates a new organisation locally and upload it to the keyserver.

This assumes you have a BitMaelum invitation token for the specific server.`,
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenDefaultVault()

		kt, err := bmcrypto.FindKeyType(*orgKeytype)
		if err != nil {
			logrus.Fatal("incorrect key type")
		}

		handlers.CreateOrganisation(v, *orgAddr, *orgFullName, *orgValidations, kt)
	},
}

var (
	orgAddr        *string
	orgFullName    *string
	orgValidations *[]string
	orgKeytype     *string
)

func init() {
	organisationCmd.AddCommand(organisationCreateCmd)

	orgAddr = organisationCreateCmd.Flags().StringP("organisation", "o", "", "Organisation address (...@<name>! part)")
	orgFullName = organisationCreateCmd.Flags().StringP("name", "n", "", "Actual name (Acme Inc.)")
	orgValidations = organisationCreateCmd.Flags().StringArray("val", nil, "validations for the organisation")
	orgKeytype = organisationCreateCmd.Flags().StringP("keytype", "k", "ed25519", "Key type to use (defaults to ED25519)")

	_ = organisationCreateCmd.MarkFlagRequired("organisation")
}
