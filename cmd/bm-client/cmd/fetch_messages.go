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
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var fetchMessagesCmd = &cobra.Command{
	Use:     "fetch-messages",
	Aliases: []string{"fetch"},
	Short:   "Retrieves messages from your account(s)",
	Long:    `Connects to the BitMaelum servers and fetches new emails that are not available on your local system.`,
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenVault()

		info := vault.GetAccountOrDefault(v, *fmAccount)
		if info == nil {
			logrus.Fatal("No account found in vault")
			os.Exit(1)
		}

		// Fetch routing info
		resolver := container.GetResolveService()
		routingInfo, err := resolver.ResolveRouting(info.RoutingID)
		if err != nil {
			logrus.Fatal("Cannot find routing ID for this account")
			os.Exit(1)
		}

		handlers.FetchMessages(info, routingInfo, *fmBox)
	},
}

var (
	// fmCheckOnly *bool
	fmAccount *string
	fmBox     *string
)

func init() {
	rootCmd.AddCommand(fetchMessagesCmd)

	fmAccount = fetchMessagesCmd.PersistentFlags().StringP("account", "a", "", "Account")
	fmBox = fetchMessagesCmd.PersistentFlags().StringP("box", "b", "", "Box to fetch")
	// fmCheckOnly = fetchMessagesCmd.Flags().Bool("check-only", false, "Check only, don't download")
}
