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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/internal/console"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/spf13/cobra"
)

var vaultImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Imports accounts or organisations into the vault",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenDefaultVault()

		_, ok := os.Stat(*viImportPath)
		if ok != nil {
			fmt.Println("import file does not exist")
			os.Exit(1)
		}

		var err error
		if *viImportPassword == "" {
			*viImportPassword, _ = console.AskPasswordPrompt("Please type your import password: ")
		}

		data, err := ioutil.ReadFile(*viImportPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		container := &vault.EncryptedContainer{}
		err = json.Unmarshal(data, &container)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if container.Type != "vault-export" {
			fmt.Println("not an vault export file")
			os.Exit(1)
		}

		b, err := vault.DecryptContainer(container, *viImportPassword)
		ex := &ExportType{}
		err = json.Unmarshal(b, &ex)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if ex.Account != nil {
			if v.HasAccount(*ex.Account.Address) {
				fmt.Println("Account already exists in vault. Refusing to import")
				os.Exit(1)
			}

			v.AddAccount(*ex.Account)
			err := v.Persist()
			if err != nil {
				fmt.Println("cannot save vault: ", err)
				os.Exit(1)
			}

			fmt.Printf("Successfully imported account: %s\n", ex.Account.Address.String())
		}

		if ex.Organisation != nil {
			if v.HasOrganisation(hash.New(ex.Organisation.Addr)) {
				fmt.Println("Organisation already exists in vault. Refusing to import")
				os.Exit(1)
			}

			v.AddOrganisation(*ex.Organisation)
			err := v.Persist()
			if err != nil {
				fmt.Println("cannot save vault: ", err)
				os.Exit(1)
			}

			fmt.Printf("Successfully imported organisation: %s\n", ex.Organisation.Addr)
		}
	},
}

var (
	viImportPassword *string
	viImportPath     *string
)

func init() {
	vaultCmd.AddCommand(vaultImportCmd)

	viImportPassword = vaultImportCmd.Flags().String("import-password", "", "password to encrypt the import file")
	viImportPath = vaultImportCmd.Flags().String("import-file", "", "path to import file")

	_ = vaultImportCmd.MarkFlagRequired("import-file")
}
