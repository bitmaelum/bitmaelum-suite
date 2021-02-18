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
	"encoding/json"
	"fmt"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/signature"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/spf13/cobra"
)

type jsonOut map[string]interface{}

// inviteCmd represents the invite command
var inviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "Invite a new user onto your server",
	Long: `This command will generate an invitation token that must be used for registering an account on your 
server. Only the specified address can register the account`,
	Run: func(cmd *cobra.Command, args []string) {
		s, _ := cmd.Flags().GetString("account")
		d, _ := cmd.Flags().GetString("duration")
		asJSON, _ := cmd.Flags().GetBool("json")

		addr, err := address.NewAddress(s)
		if err != nil {
			outError("incorrect address specified", asJSON)
			return
		}

		duration, err := internal.ValidDuration(d)
		if err != nil {
			outError("incorrect duration specified", asJSON)
			return
		}

		if config.Routing.RoutingID == "" {
			outError("cannot find routing ID for server", asJSON)
			return
		}

		validUntil := internal.TimeNow().Add(duration)
		token, err := signature.NewInviteToken(addr.Hash(), config.Routing.RoutingID, validUntil, config.Routing.KeyPair.PrivKey)
		if err != nil {
			msg := fmt.Sprintf("error while inviting address: %s", err)
			outError(msg, asJSON)
			return
		}

		if asJSON {
			output := jsonOut{
				"account": addr.String(),
				"token":   token.String(),
				"expires": validUntil.Unix(),
			}
			out, _ := json.Marshal(output)
			fmt.Printf("%s", out)
		} else {
			fmt.Printf("'%s' is allowed to register on our server until %s.\n", addr.String(), internal.TimeNow().Add(duration))
			fmt.Printf("The invitation token is: %s\n", token)
		}
	},
}

func outError(msg string, asJSON bool) {
	if !asJSON {
		fmt.Print(msg)
		return
	}

	out, _ := json.Marshal(jsonOut{"error": msg})
	fmt.Printf("%s", out)
}

func init() {
	rootCmd.AddCommand(inviteCmd)

	inviteCmd.Flags().String("account", "", "Address to register")
	inviteCmd.Flags().String("duration", "30", "NUmber of days (or duration like 1w2d3h4m6s) allowed for registration")
	inviteCmd.Flags().Bool("json", false, "Return JSON response when set")

	_ = inviteCmd.MarkFlagRequired("account")
}
