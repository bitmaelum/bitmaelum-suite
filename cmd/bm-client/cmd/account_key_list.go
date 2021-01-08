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

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var accountKeyListCmd = &cobra.Command{
	Use:   "list",
	Short: "Display account keys",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenDefaultVault()

		info, err := vault.GetAccount(v, *akAccount)
		if err != nil {
			fmt.Println("cannot find account in vault")
			os.Exit(1)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Created at", "Type", "Active", "Fingerprint"})

		for _, key := range info.Keys {
			active := ""
			if key.Active {
				active = "*"
			}

			table.Append([]string{
				internal.TimeNow().String(),
				key.PubKey.Type.String(),
				active,
				key.FingerPrint,
			})
		}

		table.Render()
	},
}

func init() {
	accountKeyCmd.AddCommand(accountKeyListCmd)
}
