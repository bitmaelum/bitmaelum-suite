package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var fetchMessagesCmd = &cobra.Command{
	Use:   "fetch-messages",
	Aliases: []string{"fetch"},
	Short: "Retrieves messages from your account(s)",
	Long: `Connects to the BitMaelum servers and fetches new emails that are not available on your local system.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("fetchMessages called witch %s and %s\n", *checkOnly, strings.Join(*addresses, " "))
	},
}

var checkOnly *bool
var addresses *[]string

func init() {
	rootCmd.AddCommand(fetchMessagesCmd)

	addresses = fetchMessagesCmd.PersistentFlags().StringArrayP("address", "a", []string{}, "Address(es) to fetch")
	checkOnly = fetchMessagesCmd.Flags().Bool("check-only", false, "Check only, don't download")
}
