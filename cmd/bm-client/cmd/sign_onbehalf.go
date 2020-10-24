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
	"fmt"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/handlers"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var onbehalfCmd = &cobra.Command{
	Use:     "sign-onbehalf",
	Short:   "Sign a key that allows mail on behalf of you",
	Long: `By signing anohter users key, you can allow it to send messages for you. This way, you can let other people 
or tools send messages without exposing your own private key.
`,
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenVault()

		info := vault.GetAccountOrDefault(v, *rAccount)
		if info == nil {
			logrus.Fatal("No account found in vault")
			os.Exit(1)
		}

		signedTargetKey, err := handlers.SignOnbehalf(info, *oTargetKey)
		if err != nil {
			logrus.Fatal("Error while signing key: ", err)
			os.Exit(1)
		}

		fmt.Printf("Signed key: %s\n", signedTargetKey)
	},
}

var (
	oAccount   *string
	oTargetKey *string
)

func init() {
	rootCmd.AddCommand(onbehalfCmd)

	oAccount = onbehalfCmd.PersistentFlags().StringP("account", "a", "", "Account")
	oTargetKey = onbehalfCmd.PersistentFlags().StringP("key", "k", "", "Key to sign")

	_ = onbehalfCmd.MarkPersistentFlagRequired("account")
	_ = onbehalfCmd.MarkPersistentFlagRequired("key")
}
