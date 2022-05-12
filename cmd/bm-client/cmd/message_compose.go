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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/handlers"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/spf13/cobra"
)

var composeCmd = &cobra.Command{
	Use:     "compose",
	Aliases: []string{"write", "send"},
	Short:   "Compose a new message",
	Long:    `This command will allow you to compose a new message and send it through your BitMaelum server`,
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenDefaultVault()

		fromInfo, err := vault.GetAccount(v, *from)
		if err != nil {
			fatal("cannot find account in vault")
		}

		toAddr, err := address.NewAddress(*to)
		if err != nil {
			fatal("cannot parse receiver address: ", err)
		}

		if *msg != "" && len(*blocks) > 0 {
			fatal("cannot specify both a messages (-m) and blocks (-b)")
		}

		// Set default message if specified
		if *msg != "" {
			*blocks = append(*blocks, "default,"+*msg)
		}

		// If no blocks are specified, and no message we assume reading a single block from stdin
		if len(*blocks) == 0 {
			var block string

			// Check if we have set $EDITOR so we can use this as our editor
			if hasEditorConfigured() {
				block, err = readFromRegularEditor()
			} else {
				// fall back to stdEditor
				block, err = readFromStdinEditor()
			}
			if err != nil {
				fatal("error reading message", err)
			}
			if len(block) == 0 {
				warn("empty message body")
			} else {
				*blocks = append(*blocks, "default,"+block)
			}
		}

		// fmt.Printf("Composing message:\n")
		// fmt.Printf("  From:    %s (%s)\n", fromInfo.Name, fromInfo.Address)
		// fmt.Printf("  To:      %s\n", *to)
		// fmt.Printf("  Subject: %s\n", *subject)
		// for i, block := range *blocks {
		// 	fmt.Printf("  Block  #%d %s\n", i, block)
		// }
		// for i, attachment := range *attachments {
		// 	fmt.Printf("  Att.   #%d %s\n", i, attachment)
		// }

		// Resolve all stuff
		resolver := container.Instance.GetResolveService()
		routingInfo, err := resolver.ResolveRouting(fromInfo.RoutingID)
		if err != nil {
			fatal("cannot find routing ID for this account")
		}

		recipientInfo, err := resolver.ResolveAddress(toAddr.Hash())
		if err != nil {
			fatal("cannot resolve recipient. Are you sure you used the correct mail address?")
		}

		// Setup addressing and compose the message
		addressing := message.NewAddressing(message.SignedByTypeOrigin)
		addressing.AddSender(fromInfo.Address, nil, fromInfo.Name, fromInfo.GetActiveKey().PrivKey, routingInfo.Routing)
		addressing.AddRecipient(toAddr, nil, &recipientInfo.PublicKey)

		err = handlers.ComposeMessage(addressing, *subject, *blocks, *attachments)
		if err != nil {
			fatal("cannot compose message: %v", err)
		}

		fmt.Println("message send successfully")
	},
}

func hasEditorConfigured() bool {
	_, err := internal.FindEditor(config.Client.Composer.Editor)
	return err == nil
}

func readFromRegularEditor() (string, error) {
	p, err := internal.FindEditor(config.Client.Composer.Editor)
	if err != nil {
		return "", err
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "bm-")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	c := exec.Command(p, tmpFile.Name())
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	err = c.Run()
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadFile(tmpFile.Name())
	return string(data), err
}

func readFromStdinEditor() (string, error) {
	if runtime.GOOS == "windows" {
		fmt.Print("\U00002709 Enter your message and press 'CTRL-Z <enter>' when done.\n")
	} else {
		fmt.Print("\U00002709 Enter your message and press 'CTRL-D' when done.\n")
	}

	data, err := ioutil.ReadAll(os.Stdin)
	return string(data), err
}

var (
	msg, from, to, subject *string
	blocks, attachments    *[]string
)

func init() {
	messageCmd.AddCommand(composeCmd)

	from = composeCmd.Flags().StringP("from", "f", "", "Sender address")
	to = composeCmd.Flags().StringP("to", "t", "", "Recipient address")
	subject = composeCmd.Flags().StringP("subject", "s", "", "Subject of the message")
	blocks = composeCmd.Flags().StringArrayP("blocks", "b", []string{}, "Message blocks")
	attachments = composeCmd.Flags().StringArrayP("attachment", "a", []string{}, "Attachments")
	msg = composeCmd.Flags().StringP("message", "m", "", "Message to send")

	_ = composeCmd.MarkFlagRequired("from")
	_ = composeCmd.MarkFlagRequired("to")
	_ = composeCmd.MarkFlagRequired("subject")
}
