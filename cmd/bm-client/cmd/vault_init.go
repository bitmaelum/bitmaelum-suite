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

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/console"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/spf13/cobra"
)

var vaultInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a new vault with a password",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		// Get either given path, or the default configuration vault path
		vaultPath := *vpath
		if vaultPath == "" {
			vaultPath = config.Client.Vault.Path
		}

		// Check if vault/path exists
		if vault.Exists(vaultPath) {
			// Vault already exists. Not creating a new vault
			fmt.Println("Vault already exists. Use --path to create a new vault on a specific path.")
			os.Exit(1)
		}

		// Check if password is given. If not, ask for password
		if vault.VaultPassword == "" {
			vault.VaultPassword, _ = console.AskDoublePassword()
		}

		// Create a new vault
		_, err := vault.Create(vaultPath, vault.VaultPassword)
		if err != nil {
			fmt.Println("error: cannot create vault: ", err)
			os.Exit(1)
		}

		fmt.Println("successfully created a new vault at ", vaultPath)
	},
}

var (
	vpath *string
)

func init() {
	vaultCmd.AddCommand(vaultInitCmd)

	vpath = vaultInitCmd.Flags().String("path", "", "Path to vault")
}
