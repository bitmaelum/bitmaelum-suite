package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/apikey"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/parse"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/spf13/cobra"
)

// apiKeyCreateCmd represents the apiKey command
var apiKeyCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create an (admin) management key for remote management and tooling",
	Example: "  apikeys --perms apikeys,invite --valid 3d --addr <hash> --desc 'My api key'",
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

		var h *hash.Hash
		if *mgAddrHash != "" {
			var err error
			h, err = hash.NewFromHash(*mgAddrHash)
			if err != nil {
				fmt.Printf("Error: incorrect hash.\n")
				os.Exit(1)
			}
		}

		// Our custom parser allows (and defaults) to using days
		validDuration, err := parse.ValidDuration(*mgValid)
		if err != nil {
			fmt.Printf("Error: incorrect duration specified.\n")
			os.Exit(1)
		}

		err = parse.Permissions(*mgPerms)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

		var key apikey.KeyType
		if *mgAdmin {
			fmt.Printf("Creating new admin key\n")
			if len(*mgPerms) > 0 {
				fmt.Printf("Error: cannot specify permissions when you create an admin key (all permissions are automatically granted)\n")
				os.Exit(1)
			}
			key = apikey.NewAdminKey(validDuration, *mgDesc)
		} else {
			fmt.Printf("Creating new regular key\n")
			if len(*mgPerms) == 0 {
				fmt.Printf("Error: need a set of permissions when generating a regular key\n")
				os.Exit(1)
			}
			key = apikey.NewKey(*mgPerms, validDuration, h, *mgDesc)
		}

		// Store API key into persistent storage
		repo := container.GetAPIKeyRepo()
		err = repo.Store(key)
		if err != nil {
			fmt.Printf("Error: cannot store key: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Your API key: %s\n", key.ID)
		if !key.ValidUntil.IsZero() {
			fmt.Printf("Key is valid until %s\n", key.ValidUntil.Format(time.RFC822))
		}
	},
}

var (
	mgAdmin    *bool
	mgPerms    *[]string
	mgValid    *string
	mgAddrHash *string
	mgDesc     *string
)

func init() {
	apiKeyCmd.AddCommand(apiKeyCreateCmd)

	mgAdmin = apiKeyCreateCmd.Flags().Bool("admin", false, "Admin key")
	mgPerms = apiKeyCreateCmd.Flags().StringSlice("permissions", []string{}, "List of permissions")
	mgValid = apiKeyCreateCmd.Flags().String("valid", "", "Days (or duration) the key is valid. Accepts 10d, or even 1h30m50s")
	mgAddrHash = apiKeyCreateCmd.Flags().String("addr", "", "Account hash for this specific api key")
	mgDesc = apiKeyCreateCmd.Flags().String("desc", "", "Description of this key")
}
