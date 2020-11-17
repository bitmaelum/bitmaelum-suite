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
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/spf13/cobra"
)

var createAccountCmd = &cobra.Command{
	Use:   "create-account",
	Short: "Create a new account",
	Long: `Create a new account locally and upload it to a BitMaelum server.

This assumes you have a BitMaelum invitation token for the specific server.`,
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenVault()

		handlers.CreateAccount(v, *addr, *name, *token, bmcrypto.KeyType(*keytype))
	},
}

var addr, name, token, keytype *string

func init() {
	rootCmd.AddCommand(createAccountCmd)

	addr = createAccountCmd.Flags().String("address", "", "Address to create")
	name = createAccountCmd.Flags().String("name", "", "Your full name")
	token = createAccountCmd.Flags().String("token", "", "Invitation token from server")
	keytype = createAccountCmd.Flags().String("keytype", "ed25519", "Type of key you want to generate (defaults to ed25519)")

	_ = createAccountCmd.MarkFlagRequired("address")
	_ = createAccountCmd.MarkFlagRequired("name")
	_ = createAccountCmd.MarkFlagRequired("token")
}
