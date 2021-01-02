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
	"io/ioutil"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-json/internal"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-json/internal/output"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/spf13/cobra"
)

var messageCmd = &cobra.Command{
	Use:   "message",
	Short: "Returns message details",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Get generic structs
		_, info, client, err := internal.GetClientAndInfo(*messageAccount)
		if err != nil {
			output.JSONErrorOut(err)
			os.Exit(1)
		}

		msg, err := client.GetMessage(info.Address.Hash(), *messageID)
		if err != nil {
			output.JSONErrorOut(err)
			os.Exit(1)
		}

		em := message.EncryptedMessage{
			ID:      msg.ID,
			Header:  &msg.Header,
			Catalog: msg.Catalog,

			GenerateBlockReader:      client.GenerateAPIBlockReader(info.Address.Hash()),
			GenerateAttachmentReader: client.GenerateAPIAttachmentReader(info.Address.Hash()),
		}

		dMsg, err := em.Decrypt(info.PrivKey)
		if err != nil {
			output.JSONErrorStrOut("cannot decrypt message")
			os.Exit(1)
		}

		out := output.JSONT{
			"header":  dMsg.Header,
			"catalog": dMsg.Catalog,
			"blocks":  []output.JSONT{},
		}

		for idx, b := range dMsg.Catalog.Blocks {
			buf, err := ioutil.ReadAll(dMsg.Catalog.Blocks[idx].Reader)
			if err != nil {
				continue
			}

			out["blocks"] = append(out["blocks"].([]output.JSONT), output.JSONT{
				"info": output.JSONT{
					"id":          b.ID,
					"compression": b.Compression,
					"checksum":    b.Checksum,
					"encoding":    b.Encoding,
					"size":        b.Size,
					"type":        b.Type,
				},
				"content": string(buf),
			})
		}

		output.JSONOut(out)
	},
}

var (
	messageAccount *string
	messageID      *string
)

func init() {
	rootCmd.AddCommand(messageCmd)

	messageAccount = messageCmd.Flags().StringP("account", "a", "", "Account")
	messageID = messageCmd.Flags().String("id", "", "Message ID")
	_ = messageCmd.MarkFlagRequired("account")
	_ = messageCmd.MarkFlagRequired("message")
}
