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
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/key"
	"github.com/bitmaelum/bitmaelum-suite/internal/parse"
	"github.com/spf13/cobra"
)

// apiKeyCreateCmd represents the apiKey command
var apiKeyCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create an (admin) management key for remote management and tooling",
	Example: "  apikeys --perms apikeys,invite --valid 3d --desc 'My api key'",
	Long: `This command will generate an management key that can be used to administer commands through the HTTPS server. By default this is disabled, 
but can be enabled with the server.management.enabled flag in the server configuration file.

  Management permissions:
    flush            Enables remote flushing of all queues so mail is processed immediately.
    mail             Allows sending mail without a registered account.
    invite           Generate invites remotely.
    apikeys          Remove or add API keys (except admin keys).

  Account permissions:
    get-headers      Retrieve message headers from a specific account or accounts.

Note: Creating an admin key can only be done locally on the mail-server.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if !config.Server.Management.Enabled {
			fmt.Printf("Warning: remote management is not enabled on this server. You need to enable this in your configuration first.\n\n")
		}

		// Our custom parser allows (and defaults) to using days
		validDuration, err := parse.ValidDuration(*mgValid)
		if err != nil {
			fmt.Printf("Error: incorrect duration specified.\n")
			os.Exit(1)
		}

		var expires = time.Time{}
		if validDuration > 0 {
			expires = time.Now().Add(validDuration)
		}

		err = parse.ManagementPermissions(*mgPerms)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

		var k key.APIKeyType
		if *mgAdmin {
			fmt.Printf("Creating new admin key\n")
			if len(*mgPerms) > 0 {
				fmt.Printf("Error: cannot specify permissions when you create an admin key (all permissions are automatically granted)\n")
				os.Exit(1)
			}
			k = key.NewAPIAdminKey(expires, *mgDesc)
		} else {
			fmt.Printf("Creating new regular key\n")
			if len(*mgPerms) == 0 {
				fmt.Printf("Error: need a set of permissions when generating a regular key\n")
				os.Exit(1)
			}
			k = key.NewAPIKey(*mgPerms, expires, *mgDesc)
		}

		// Store API key into persistent storage
		repo := container.GetAPIKeyRepo()
		err = repo.Store(k)
		if err != nil {
			fmt.Printf("Error: cannot store key: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Your API key: %s\n", k.ID)
		if !k.Expires.IsZero() {
			fmt.Printf("Key is valid until %s\n", k.Expires.Format(time.RFC822))
		}
	},
}

var (
	mgAdmin *bool
	mgPerms *[]string
	mgValid *string
	mgDesc  *string
)

func init() {
	apiKeyCmd.AddCommand(apiKeyCreateCmd)

	mgAdmin = apiKeyCreateCmd.Flags().Bool("admin", false, "Admin key")
	mgPerms = apiKeyCreateCmd.Flags().StringSlice("permissions", []string{}, "List of permissions")
	mgValid = apiKeyCreateCmd.Flags().String("valid", "", "Days (or duration) the key is valid. Accepts 10d, or even 1h30m50s")
	mgDesc = apiKeyCreateCmd.Flags().String("desc", "", "Description of this key")
}
