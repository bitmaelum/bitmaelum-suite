package cmd

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-config/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/apikey"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/spf13/cobra"
	"os"
	"time"
)

// apiKeyCmd represents the apiKey command
var apiKeyCmd = &cobra.Command{
	Use:     "api-key",
	Aliases: []string{"apikey", "key"},
	Short:   "Create an (admin) management key for remote management",
	Example: "  apikeys --perms apikeys,invite --valid 3d",
	Long: `This command will generate an management key that can be used to administer commands through the HTTPS server. By default this is disabled, 
but can be enabled with the server.management.enabled flag in the server configuration file.

Permission list:
    
    flush            Enables remote flushing of all queues so mail is processed immediately
    mail             Allows sending mail without a registered account
    invite           Generate invites remotely
    apikeys          Remove or add API keys (except admin keys)

Creating an admin key can only be done locally.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if !config.Server.Management.Enabled {
			fmt.Printf("Warning: remote management is not enabled on this server. You need to enable this in your configuration first.\n\n")
		}

		// Our custom parser allows (and defaults) to using days
		validDuration, err := internal.ParseValidDuration(*mgValid)
		if err != nil {
			fmt.Printf("Error: incorrect duration specified")
			os.Exit(1)
		}

		err = internal.ParsePermissions(*mgPerms)
		if err != nil {
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}

		var key apikey.KeyType
		if *mgAdmin {
			fmt.Printf("Creating new admin key\n")
			if len(*mgPerms) > 0 {
				fmt.Printf("Error: cannot specify permissions when you create an admin key (all permissions are automatically granted)")
				os.Exit(1)
			}
			key = apikey.NewAdminKey(validDuration)
		} else {
			fmt.Printf("Creating new regular key\n")
			if len(*mgPerms) == 0 {
				fmt.Printf("Error: need a set of permissions when generating a regular key")
				os.Exit(1)
			}
			key = apikey.NewKey(*mgPerms, validDuration)
		}

		// Store API key into persistent storage
		repo := container.GetAPIKeyRepo()
		err = repo.Store(key)
		if err != nil {
			fmt.Printf("Error: cannot store key: %s", err)
			os.Exit(1)
		}

		fmt.Printf("Your API key: %s\n", key.ID)
		if !key.ValidUntil.IsZero() {
			fmt.Printf("Key is valid until %s\n", key.ValidUntil.Format(time.RFC822))
		}
	},
}

var (
	mgAdmin *bool
	mgPerms *[]string
	mgValid *string
)

func init() {
	rootCmd.AddCommand(apiKeyCmd)

	mgAdmin = apiKeyCmd.Flags().Bool("admin", false, "Admin key")
	mgPerms = apiKeyCmd.Flags().StringSlice("permissions", []string{}, "List of permissions")
	mgValid = apiKeyCmd.Flags().String("valid", "", "Days (or duration) the key is valid. Accepts 10d, or even 1h30m50s")
}
