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
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/handlers"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal"
	pkginternal "github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/spf13/cobra"
)

var messageListCmd = &cobra.Command{
	Use:   "list",
	Short: "Displays a list of messages from your account(s)",
	Long:  `Retrieves and displays a list of message found on your remote server`,
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenDefaultVault()

		var since time.Time

		if *lmNew && *lmSince != "" {
			fmt.Println("You can specify either --new or --since, but not both")
			os.Exit(1)
		}

		if *lmSince != "" {
			d, err := pkginternal.ParseDuration(*lmSince)
			if err != nil {
				fmt.Println("incorrect --since format. Use the following format: 1y3w4d5h13m")
				os.Exit(1)
			}
			since = pkginternal.TimeNow().Add(-1 * d)
		}

		if *lmNew {
			since = internal.GetReadTime()
		}


		// Get account or all accounts
		accounts := v.Store.Accounts
		if *lmAccount != "" {
			acc, err := vault.GetAccount(v, *lmAccount)
			if err == nil {
				accounts = []vault.AccountInfo{*acc}
			}
		}


		msgCount := handlers.ListMessages(accounts, since)
		if msgCount == 0 {
			if *lmNew {
				fmt.Println("* No new messages found")
			} else {
				fmt.Println("* No messages since ", since.Format(time.RFC822))
			}
		}

		internal.SaveReadTime(pkginternal.TimeNow())
	},
}

var (
	lmNew     *bool
	lmAccount *string
	lmSince   *string
)

func init() {
	messageCmd.AddCommand(messageListCmd)

	lmAccount = messageListCmd.Flags().StringP("account", "a", "", "Account")

	lmNew = messageListCmd.Flags().BoolP("new", "n", false, "Display new messages only")
	lmSince = messageListCmd.Flags().StringP("since", "s", "", "Display messages since the specific duration (accepts 1y1w1d1h)")
}
