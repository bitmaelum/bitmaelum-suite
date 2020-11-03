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

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/handlers"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var otpCmd = &cobra.Command{
	Use:   "otp",
	Short: "Generate a OTP to authenticate to a server",
	Long:  `It will generate a OTP code that can be used to authenticate to a third party server`,
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenVault()

		accountToUse := vault.GetAccountOrDefault(v, *otpAccount)
		if accountToUse == nil {
			logrus.Fatal("No account found in vault")
			os.Exit(1)
		}

		handlers.OtpGenerate(accountToUse, otpServer)

	},
}

var (
	otpServer  *string
	otpAccount *string
)

func init() {
	rootCmd.AddCommand(otpCmd)

	otpAccount = otpCmd.Flags().String("account", "a", "Account to be used for authentication")
	otpServer = otpCmd.Flags().String("server", "s", "Server domain name to get public key from")

	_ = otpCmd.MarkFlagRequired("server")
}
