package cmd

import (
	"github.com/bitmaelum/bitmaelum-server/bm-client/handlers"
	"github.com/spf13/cobra"
)

var createAccountCmd = &cobra.Command{
	Use:   "create-account",
	Short: "Create a new account",
	Long: `Create a new account locally and upload it to a BitMaelum servrer.

This assumes you have a BitMaelum invitation token for the specific server.`,
	Run: func(cmd *cobra.Command, args []string) {
		pwd, err := cmd.Flags().GetString("password")
		if err != nil {
			panic(err)
		}

		handlers.CreateAccount(*address, *name, *organisation, *server, []byte(pwd))
	},
}

var address, name, organisation, server *string

func init() {
	rootCmd.AddCommand(createAccountCmd)

	address = createAccountCmd.Flags().String("address", "", "Address to create")
	name = createAccountCmd.Flags().String("name", "", "Name")
	organisation = createAccountCmd.Flags().String("org", "", "Organisation")
	server = createAccountCmd.Flags().String("server", "", "Server to store the account")

	_ = createAccountCmd.MarkFlagRequired("address")
	_ = createAccountCmd.MarkFlagRequired("name")
	_ = createAccountCmd.MarkFlagRequired("server")
}
