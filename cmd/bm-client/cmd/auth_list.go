// Copyright (c) 2022 BitMaelum Authors
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

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var authListCmd = &cobra.Command{
	Use:   "list",
	Short: "Display all authorized keys",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		// Get generic structs
		_, info, client, err := internal.GetClientAndInfo(*authAccount)
		if err != nil {
			logrus.Fatal(err)
			os.Exit(1)
		}

		keys, err := client.ListAuthKeys(info.Address.Hash())
		if err != nil {
			logrus.Fatal(err)
			os.Exit(1)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Fingerprint", "Valid until", "Description"})

		for _, key := range keys {
			// don't display zero times
			expiry := key.Expires.Format(time.ANSIC)
			if key.Expires.Unix() == 0 {
				expiry = ""
			}

			table.Append([]string{
				key.Fingerprint,
				expiry,
				key.Description,
			})
		}

		table.Render()
	},
}

func init() {
	authCmd.AddCommand(authListCmd)
}
