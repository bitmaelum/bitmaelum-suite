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
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/internal/console"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/spf13/cobra"
)

// ExportType is the structure that holds either an account or an organsation info
type ExportType struct {
	Type         string
	Account      *vault.AccountInfo
	Organisation *vault.OrganisationInfo
}

var vaultExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Exports accounts or organisations from the vault",
	Run: func(cmd *cobra.Command, args []string) {
		if *veOrganisation != "" && *veAccount != "" {
			fmt.Println("You can only specify either --organisation or --account. Not both")
			os.Exit(1)
		}
		if *veOrganisation == "" && *veAccount == "" {
			fmt.Println("You must specify either --organisation or --account")
			os.Exit(1)
		}

		v := vault.OpenDefaultVault()

		var (
			export *ExportType
			err    error
		)

		if *veAccount != "" {
			export, err = exportAccount(v, *veAccount)
		}
		if *veOrganisation != "" {
			export, err = exportOrganisation(v, *veOrganisation)
		}

		if *veExportPassword == "" {
			*veExportPassword, err = console.AskDoublePasswordPrompt("Please type your export password: ", "Retry your export password: ")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		b, err := vault.EncryptContainer(export, *veExportPassword, "vault-export", 1)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = ioutil.WriteFile(*veExportPath, b, 0600)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("successfully created export file")
	},
}

func exportAccount(v *vault.Vault, account string) (*ExportType, error) {
	addr, err := address.NewAddress(account)
	if err != nil {
		return nil, errors.New("incorrect account specified")
	}

	info, err := v.GetAccountInfo(*addr)
	if err != nil {
		fmt.Println("Account not found")
		os.Exit(1)
	}

	return &ExportType{
		Type:    "account",
		Account: info,
	}, nil
}

func exportOrganisation(v *vault.Vault, org string) (*ExportType, error) {
	orgHash := hash.New(org)

	info, err := v.GetOrganisationInfo(orgHash)
	if err != nil {
		fmt.Println("Organisation not found")
		os.Exit(1)
	}

	return &ExportType{
		Type:         "org",
		Organisation: info,
	}, nil
}

var (
	veAccount        *string
	veOrganisation   *string
	veExportPassword *string
	veExportPath     *string
)

func init() {
	vaultCmd.AddCommand(vaultExportCmd)

	veOrganisation = vaultExportCmd.Flags().StringP("organisation", "o", "", "org name to export")
	veAccount = vaultExportCmd.Flags().StringP("account", "a", "", "account to export")
	veExportPassword = vaultExportCmd.Flags().String("export-password", "", "password to encrypt the export file")
	veExportPath = vaultExportCmd.Flags().String("export-file", "", "path to export file")

	_ = vaultExportCmd.MarkFlagRequired("export-file")
}
