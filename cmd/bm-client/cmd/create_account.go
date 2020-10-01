package cmd

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/handlers"
	"github.com/spf13/cobra"
)

var createAccountCmd = &cobra.Command{
	Use:   "create-account",
	Short: "Create a new account",
	Long: `Create a new account locally and upload it to a BitMaelum server.

This assumes you have a BitMaelum invitation token for the specific server.`,
	Run: func(cmd *cobra.Command, args []string) {
		vault := OpenVault()

		handlers.CreateAccount(vault, *addr, *name, *token)
	},
}

var addr, name, token *string

func init() {
	rootCmd.AddCommand(createAccountCmd)

	addr = createAccountCmd.Flags().String("address", "", "Address to create")
	name = createAccountCmd.Flags().String("name", "", "Your full name")
	token = createAccountCmd.Flags().String("token", "", "Invitation token from server")

	_ = createAccountCmd.MarkFlagRequired("address")
	_ = createAccountCmd.MarkFlagRequired("name")
	_ = createAccountCmd.MarkFlagRequired("token")
}
