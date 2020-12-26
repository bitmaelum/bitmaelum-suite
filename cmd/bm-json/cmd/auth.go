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

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-json/internal"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-json/internal/output"
	"github.com/spf13/cobra"
)

var authKeyCmd = &cobra.Command{
	Use:   "auth",
	Short: "Returns auth key info",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Get generic structs
		_, info, client, err := internal.GetClientAndInfo(*authAccount)
		if err != nil {
			output.JSONErrorOut(err)
			os.Exit(1)
		}

		authKeys, err := client.ListAuthKeys(info.Address.Hash())
		if err != nil {
			output.JSONErrorOut(err)
			os.Exit(1)
		}

		out := []output.JSONT{}
		for _, ak := range authKeys {

			// don't display zero times
			expiry := ak.Expires.Format(time.ANSIC)
			if ak.Expires.Unix() == 0 {
				expiry = ""
			}

			out = append(out, output.JSONT{
				"id":           ak.Fingerprint,
				"expires":      expiry,
				"public_key":   ak.PublicKey,
				"description":  ak.Description,
				"address_hash": ak.AddressHash,
				"signature":    ak.Signature,
			})
		}

		output.JSONOut(out)
	},
}

var authAccount *string

func init() {
	rootCmd.AddCommand(authKeyCmd)

	authAccount = authKeyCmd.Flags().StringP("account", "a", "", "Account")
	_ = authKeyCmd.MarkFlagRequired("account")
}
