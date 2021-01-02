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
	"github.com/spf13/cobra"
)

var webhookCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new webhook",
}

var whEvent *string

func init() {
	webhookCmd.AddCommand(webhookCreateCmd)

	webhookCreateCmd.SetHelpTemplate(`The following webhooks events (--event or -e) are supported:

  localdelivery       When a message is delivered locally
  faileddelivery      When a message cannot be delivered
  remotedelivery      When a message is delivered remotely (outgoing mail)

  apikeycreated       When an API key has been created
  apikeydeleted       When an API key has been deleted
  apikeyupdated       When an API key has been updated

  authkeycreated      When an auth key has been created
  authkeydeleted      When an auth key has been deleted
  authkeyupdated      When an auth key has been updated

  webhookcreated      When a webhook has been created
  webhookdeleted      When a webhook has been deleted
  webhookupdated      When a webhook has been updated

or use "all" to trigger on ALL webhook events.
			
`)

	whEvent = webhookCreateCmd.PersistentFlags().StringP("event", "e", "", "Event to use (see 'bm-client webhook create' for options)")

	_ = webhookCreateCmd.MarkPersistentFlagRequired("event")
}
