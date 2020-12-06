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
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-json/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-json/internal/output"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/spf13/cobra"
)

var apiKeyCmd = &cobra.Command{
	Use:   "api",
	Short: "Returns api key info",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenVault()
		info := vault.GetAccountOrDefault(v, *apiAccount)
		if info == nil {
			output.JSONErrorStrOut("Account not found")
			os.Exit(1)
		}

		resolver := container.Instance.GetResolveService()
		routingInfo, err := resolver.ResolveRouting(info.RoutingID)
		if err != nil {
			output.JSONErrorStrOut("Cannot find routing ID for this account")
			os.Exit(1)
		}

		client, err := api.NewAuthenticated(*info.Address, &info.PrivKey, routingInfo.Routing, internal.JwtJSONErrorFunc)
		if err != nil {
			output.JSONErrorOut(err)
			os.Exit(1)
		}

		apiKeys, err := client.ListAPIKeys(info.Address.Hash())
		if err != nil {
			output.JSONErrorOut(err)
			os.Exit(1)
		}

		if apiKeys == nil {

		}

		var out []output.JSONT
		for _, ak := range apiKeys {

			// don't display zero times
			expiry := ak.Expires.Format(time.ANSIC)
			if ak.Expires.Unix() == 0 {
				expiry = ""
			}

			ah := ""
			if ak.AddressHash != nil {
				ah = ak.AddressHash.String()
			}

			out = append(out, output.JSONT{
				"id": ak.ID,
				"expires": expiry,
				"permissions": ak.Permissions,
				"description": ak.Desc,
				"address_hash": ah,
				"is_admin": ak.Admin,
			})
		}

		output.JSONOut(out)
	},
}

var apiAccount *string

func init() {
	rootCmd.AddCommand(apiKeyCmd)

	apiAccount = apiKeyCmd.Flags().StringP("account", "a", "", "Account")
	_ = apiKeyCmd.MarkFlagRequired("account")
}
