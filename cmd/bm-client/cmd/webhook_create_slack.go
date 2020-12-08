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
	"github.com/bitmaelum/bitmaelum-suite/internal/webhook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var webhookCreateSlackCmd = &cobra.Command{
	Use:   "slack",
	Short: "Creates a new slack webhook",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		// Validate event
		evt, err := webhook.NewEventFromString(*whEvent)
		if err != nil {
			fmt.Println("unknown event: ", *whEvent)
			fmt.Println("")

			_ = webhookCreateCmd.Help()
			os.Exit(1)
		}

		v := vault.OpenVault()
		info := vault.GetAccountOrDefault(v, *whAccount)

		resolver := container.Instance.GetResolveService()
		routingInfo, err := resolver.ResolveRouting(info.RoutingID)
		if err != nil {
			logrus.Fatal("Cannot find routing ID for this account")
			os.Exit(1)
		}

		cfg := &webhook.ConfigSlack{
			WebhookURL: *whsUrl,
			Channel:    *whsChannel,
			Username:   *whsUsername,
			IconEmoji:  *whsIconEmoji,
			IconUrl:    *whsIconUrl,
			Template:   *whsTemplate,
		}
		wh, err := webhook.NewWebhook(info.Address.Hash(), evt, webhook.TypeSlack, cfg)
		if err != nil {
			logrus.Fatal("Cannot create webhook")
			os.Exit(1)
		}

		client, err := api.NewAuthenticated(*info.Address, &info.PrivKey, routingInfo.Routing, internal.JwtErrorFunc)
		if err != nil {
			logrus.Fatal(err)
			os.Exit(1)
		}

		wh, err = client.CreateWebhook(*wh)
		if err != nil {
			logrus.Fatal(err)
			os.Exit(1)
		}

		fmt.Println("Created webhook %s", wh.ID)
	},
}

var (
	whsUrl       *string
	whsChannel   *string
	whsUsername  *string
	whsIconEmoji *string
	whsIconUrl   *string
	whsTemplate  *string
)

func init() {
	webhookCreateCmd.AddCommand(webhookCreateSlackCmd)

	whsUrl = webhookCreateSlackCmd.Flags().String("url", "", "Slack webhook URL")
	whsChannel = webhookCreateSlackCmd.Flags().String("channel", "", "Optional channel to post to")
	whsUsername = webhookCreateSlackCmd.Flags().String("username", "", "Optional username to post from")
	whsIconEmoji = webhookCreateSlackCmd.Flags().String("icon_emoji", "", "Optional bot icon emoji")
	whsIconUrl = webhookCreateSlackCmd.Flags().String("icon_url", "", "Optional bot icon url")
	whsTemplate = webhookCreateSlackCmd.Flags().String("template", "", "Optional text/template")

	_ = webhookCreateSlackCmd.MarkFlagRequired("url")
}
