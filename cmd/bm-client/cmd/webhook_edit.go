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
	"encoding/json"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal"
	internal2 "github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var webhookEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit webhook",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Get generic structs
		_, info, client, err := internal.GetClientAndInfo(*whAccount)
		if err != nil {
			logrus.Fatal(err)
			os.Exit(1)
		}

		// Get webhook and edit
		src, err := client.GetWebhook(info.Address.Hash(), *wheditID)
		if err != nil {
			logrus.Fatal("error while editing webhook: ", err)
			os.Exit(1)
		}

		// Unmarshal src-config for editing
		var srcConfig interface{}
		err = json.Unmarshal([]byte(src.Config), &srcConfig)
		if err != nil {
			logrus.Fatal("error while editing webhook: ", err)
			os.Exit(1)
		}

		// Edit into dstConig
		var dstConfig interface{}
		err = internal2.JSONFileEditor(srcConfig, &dstConfig)
		if err != nil {
			logrus.Fatal("error while editing webhook: ", err)
			os.Exit(1)
		}

		// Marshal dstconfig back into src config
		data, err := json.Marshal(dstConfig)
		if err != nil {
			logrus.Fatal("error while editing webhook: ", err)
			os.Exit(1)
		}
		src.Config = string(data)

		// And update webhook
		err = client.UpdateWebhook(info.Address.Hash(), *wheditID, *src)
		if err != nil {
			logrus.Fatal("cannot update webhook: ", err)
			os.Exit(1)
		}
	},
}

var wheditID *string

func init() {
	wheditID = webhookEditCmd.Flags().String("id", "", "webhook ID to edit")

	webhookCmd.AddCommand(webhookEditCmd)
}
