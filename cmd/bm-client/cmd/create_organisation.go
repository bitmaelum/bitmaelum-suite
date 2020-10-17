// Copyright (c) 2020 BitMaelum Authors
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
	"github.com/spf13/cobra"
)

var createOrganisationCmd = &cobra.Command{
	Use:   "create-organisation",
	Short: "Create a new organisation",
	Long: `Create a new organisation locally and upload it to the keyserver.

This assumes you have a BitMaelum invitation token for the specific server.`,
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenVault()

		handlers.CreateOrganisation(v, *orgAddr, *orgFullName, *orgValidations, *rsa)
	},
}

var (
	orgAddr        *string
	orgFullName    *string
	orgValidations *[]string
)

func init() {
	rootCmd.AddCommand(createOrganisationCmd)

	orgAddr = createOrganisationCmd.Flags().StringP("org", "o", "", "Organisation address (...@<name>! part)")
	orgFullName = createOrganisationCmd.Flags().StringP("name", "n", "", "Actual name (Acme Inc.)")
	orgValidations = createOrganisationCmd.Flags().StringArray("val", nil, "validations for the organisation")

	_ = createAccountCmd.MarkFlagRequired("org")
}
