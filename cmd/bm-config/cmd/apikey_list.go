package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// apiKeyListCmd represents the apiKey command
var apiKeyListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all current keys",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("listing keys")
	},
}

func init() {
	apiKeyCmd.AddCommand(apiKeyListCmd)

	// mgAdmin = apiKeyCreateCmd.Flags().Bool("admin", false, "Admin key")
	// mgPerms = apiKeyCreateCmd.Flags().StringSlice("permissions", []string{}, "List of permissions")
	// mgValid = apiKeyCreateCmd.Flags().String("valid", "", "Days (or duration) the key is valid. Accepts 10d, or even 1h30m50s")
	// mgAddrHash = apiKeyCreateCmd.Flags().String("addr", "", "Account hash for this specific api key")
	// mgDesc = apiKeyCreateCmd.Flags().String("desc", "", "Description of this key")
}
