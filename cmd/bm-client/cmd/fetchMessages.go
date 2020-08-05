package cmd

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/handlers"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/spf13/cobra"
)

var fetchMessagesCmd = &cobra.Command{
	Use:     "fetch-messages",
	Aliases: []string{"fetch"},
	Short:   "Retrieves messages from your account(s)",
	Long:    `Connects to the BitMaelum servers and fetches new emails that are not available on your local system.`,
	Run: func(cmd *cobra.Command, args []string) {
		addr, err := address.New(*account)
		if err != nil {
			panic(err)
		}
		info, err := Vault.GetAccountInfo(*addr)
		if err != nil {
			panic(err)
		}

		handlers.FetchMessages(info, *box, *checkOnly)
	},
}

var checkOnly *bool
var account *string
var box *string

func init() {
	rootCmd.AddCommand(fetchMessagesCmd)

	account = fetchMessagesCmd.PersistentFlags().StringP("account", "a", "", "Account")
	box = fetchMessagesCmd.PersistentFlags().StringP("box", "b", "", "Box to fetch")
	checkOnly = fetchMessagesCmd.Flags().Bool("check-only", false, "Check only, don't download")
}
