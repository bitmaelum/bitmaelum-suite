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

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var webhookListCmd = &cobra.Command{
	Use:   "list",
	Short: "Display all webhooks for this account on the server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenVault()
		info := vault.GetAccountOrDefault(v, *whAccount)

		resolver := container.Instance.GetResolveService()
		routingInfo, err := resolver.ResolveRouting(info.RoutingID)
		if err != nil {
			logrus.Fatal("Cannot find routing ID for this account")
			os.Exit(1)
		}

		client, err := api.NewAuthenticated(*info.Address, &info.PrivKey, routingInfo.Routing, internal.JwtErrorFunc)
		if err != nil {
			logrus.Fatal(err)
			os.Exit(1)
		}

		webhooks, err := client.ListWebhooks(info.Address.Hash())
		if err != nil {
			logrus.Fatal("cannot fetch webhooks: ", err)
			os.Exit(1)
		}

		if len(webhooks) == 0 {
			fmt.Println("No webhooks are found for this account")
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Event", "Type", "Enabled"})

		for _, wh := range webhooks {
			enabled := "-"
			if wh.Enabled {
				enabled = "enabled"
			}

			table.Append([]string{
				wh.ID,
				wh.Event.String(),
				wh.Type.String(),
				enabled,
			})
		}

		table.Render()
	},
}

func init() {
	webhookCmd.AddCommand(webhookListCmd)
}
