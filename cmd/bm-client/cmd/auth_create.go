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
	"os"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/handlers"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var authCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new authorized key to send messages on your behalf",
	Long: `By authorizing another users public key, you can allow it to send messages for you. This way, you can let other people or tools send messages without exposing your own private key.
`,
	Run: func(cmd *cobra.Command, args []string) {
		targetKey, err := bmcrypto.NewPubKey(*authPublicKey)
		if err != nil {
			logrus.Fatal("Incorrect public key to sign")
			os.Exit(1)
		}

		v := vault.OpenVault()
		info := vault.GetAccountOrDefault(v, *authAccount)
		if info == nil {
			logrus.Fatal("No account found in vault")
			os.Exit(1)
		}

		err = handlers.CreateAuthorizedKey(info, targetKey, *authValidUntil, *authDesc)
		if err != nil {
			logrus.Fatal(err)
			os.Exit(1)
		}

		logrus.Print("Created new authorized key")
	},
}

var (
	authPublicKey  *string
	authValidUntil *time.Duration
	authDesc       *string
)

func init() {
	authCmd.AddCommand(authCreateCmd)

	authPublicKey = authCreateCmd.Flags().StringP("key", "k", "", "Key to sign")
	authValidUntil = authCreateCmd.Flags().Duration("duration", time.Duration(0), "Time valid")
	authDesc = authCreateCmd.Flags().StringP("desc", "d", "", "Value to set")

	_ = authCreateCmd.MarkFlagRequired("desc")
	_ = authCreateCmd.MarkFlagRequired("key")
}
