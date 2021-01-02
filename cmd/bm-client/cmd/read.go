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
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal/container"
	pkginternal "github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:     "read",
	Aliases: []string{"read-message", "r"},
	Short:   "Read messages from your account",
	Long: `Read message from your account
`,
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenVault()

		info := vault.GetAccountOrDefault(v, *rAccount)
		if info == nil {
			logrus.Fatal("* No account found in vault")
			os.Exit(1)
		}

		// Fetch routing info
		resolver := container.Instance.GetResolveService()
		routingInfo, err := resolver.ResolveRouting(info.RoutingID)
		if err != nil {
			logrus.Fatal("* Cannot find routing ID for this account")
			os.Exit(1)
		}

		var since time.Time

		if *rNew && *rSince != "" {
			fmt.Println("* You can specify either --new or --since, but not both")
			os.Exit(1)
		}

		if *rSince != "" {
			d, err := pkginternal.ParseDuration(*rSince)
			if err != nil {
				fmt.Println("* Incorrect --since format. Use the following format: 1y3w4d5h13m")
				os.Exit(1)
			}
			since = pkginternal.TimeNow().Add(-1 * d)
		}

		if *rNew {
			since = internal.GetReadTime()
		}

		handlers.ReadMessages(info, routingInfo, *rBox, *rMessageID, since)

		internal.SaveReadTime(pkginternal.TimeNow())
	},
}

var (
	rAccount   *string
	rBox       *string
	rMessageID *string
	rSince     *string
	rNew       *bool
)

func init() {
	rootCmd.AddCommand(readCmd)

	rAccount = readCmd.Flags().StringP("account", "a", "", "Account")

	rBox = readCmd.Flags().StringP("box", "b", "", "Box to fetch")
	rMessageID = readCmd.Flags().String("id", "", "Message ID")
	rNew = readCmd.Flags().BoolP("new", "n", false, "Read new messages only")
	rSince = readCmd.Flags().StringP("since", "s", "", "Read messages since the specific duration (accepts 1y1w1d1h)")

	_ = readCmd.MarkFlagRequired("account")
}
