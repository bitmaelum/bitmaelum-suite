package cmd

import (
	"github.com/bitmaelum/bitmaelum-server/bm-client/core"
	"github.com/spf13/cobra"
)

var listAccountsCmd = &cobra.Command{
	Use:   "list-accounts",
	Aliases: []string{"list-account", "la", "list"},
	Short: "List your accounts",
	Long: `Displays a list of all your accounts currently available`,
	Run: func(cmd *cobra.Command, args []string) {
		core.ListAccounts()
	},
}

func init() {
	rootCmd.AddCommand(listAccountsCmd)
}
