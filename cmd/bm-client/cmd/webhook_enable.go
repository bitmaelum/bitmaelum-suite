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
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var webhookEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable webhook",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Get generic structs
		_, info, client, err := internal.GetClientAndInfo(*whAccount)
		if err != nil {
			logrus.Fatal(err)
			os.Exit(1)
		}

		err = client.EnableWebhook(info.Address.Hash(), *wheID)
		if err != nil {
			logrus.Fatal("cannot enable webhook: ", err)
			os.Exit(1)
		}

		fmt.Println("Webhook is enabled")
	},
}

var wheID *string

func init() {
	wheID = webhookEnableCmd.Flags().String("id", "", "webhook ID to enable")

	webhookCmd.AddCommand(webhookEnableCmd)
}
