package cmd

import (
	"github.com/bitmaelum/bitmaelum-server/cmd/bm-client/handlers"
	"github.com/spf13/cobra"
)

var createAccountCmd = &cobra.Command{
	Use:   "create-account",
	Short: "Create a new account",
	Long: `Create a new account locally and upload it to a BitMaelum servrer.

This assumes you have a BitMaelum invitation token for the specific server.`,
	Run: func(cmd *cobra.Command, args []string) {
		handlers.CreateAccount(&Vault, *addr, *name, *organisation, *server, *token)
	},
}

var addr, name, organisation, server, token *string

func init() {
	rootCmd.AddCommand(createAccountCmd)

	addr = createAccountCmd.Flags().String("address", "", "Address to create")
	name = createAccountCmd.Flags().String("name", "", "Your full name")
	organisation = createAccountCmd.Flags().String("org", "", "Organisation")
	server = createAccountCmd.Flags().String("server", "", "Server to store the account")
	token = createAccountCmd.Flags().String("token", "", "Invitation token from server")

	_ = createAccountCmd.MarkFlagRequired("address")
	_ = createAccountCmd.MarkFlagRequired("name")
	_ = createAccountCmd.MarkFlagRequired("server")
	_ = createAccountCmd.MarkFlagRequired("token")
}
